package handler

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/infrastructure/repository"
	"security_chat_app/internal/interface/markup"
	"security_chat_app/internal/interface/middleware"
	"security_chat_app/internal/utils/icons"
)

// ProfileData プロフィールページのデータ構造体
type ProfileData struct {
	IsLoggedIn     bool
	LoggedInUserID string
	User           *domain.User
	Posts          []domain.Post
	Replies        []domain.Post
	Likes          []domain.Post
}

// プロフィールページの表示
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// URLからユーザーIDを取得
	path := r.URL.Path
	var targetUserID string
	if path == "/profile" || path == "/profile/" {
		targetUserID = session.User.ID
	} else {
		targetUserID = path[len("/profile/"):]
		if targetUserID == "" {
			http.Error(w, "ユーザーIDが指定されていません", http.StatusBadRequest)
			return
		}
	}

	// ユーザー情報の取得
	user, err := repository.GetUserByID(targetUserID)
	if err != nil {
		http.Error(w, "ユーザー情報の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// アイコンが設定されていない場合はデフォルトアイコンを設定
	if user.Icon == "" {
		randomNum := rand.Intn(7)
		defaultIconPath := fmt.Sprintf(icons.DefaultIconPath+"/default_icon_%s.png", icons.DefaultIconNames[randomNum])
		iconURL, er := firebase.GetDefaultIconURL(defaultIconPath)
		if er != nil {
			fmt.Printf("デフォルトアイコンの取得に失敗: %v\n", err)
			http.Error(w, "デフォルトアイコンの取得に失敗しました", http.StatusInternalServerError)
			return
		}

		// ユーザーのIconURLを更新
		user.Icon = iconURL
		err = firebase.UpdateField("users", user.ID, "Icon", iconURL)
		if err != nil {
			http.Error(w, "アイコンURLの更新に失敗しました", http.StatusInternalServerError)
			return
		}
	}

	// 投稿、返信、いいねを取得
	posts, err := firebase.GetUserPosts(user.ID)
	if err != nil {
		fmt.Printf("投稿の取得に失敗: %v\n", err)
		posts = []domain.Post{}
	}

	replies, err := firebase.GetUserReplies(user.ID)
	if err != nil {
		fmt.Printf("返信の取得に失敗: %v\n", err)
		replies = []domain.Post{}
	}

	likes, err := firebase.GetUserLikedPosts(user.ID)
	if err != nil {
		fmt.Printf("いいねの取得に失敗: %v\n", err)
		likes = []domain.Post{}
	}

	// 最終更新日時を現在時刻に更新 (自分のプロフィールの場合のみ更新すべきか検討)
	if targetUserID == session.User.ID {
		user.UpdatedAt = time.Now()
		err = firebase.UpdateField("users", user.ID, "UpdatedAt", user.UpdatedAt)
		if err != nil {
			http.Error(w, "最終更新日時の更新に失敗しました", http.StatusInternalServerError)
			return
		}
	}

	// プロフィールデータの作成
	data := ProfileData{
		IsLoggedIn:     true,
		LoggedInUserID: session.User.ID,
		User:           user,
		Posts:          posts,
		Replies:        replies,
		Likes:          likes,
	}

	// テンプレートを描画
	markup.GenerateHTML(w, data, "layout", "header", "profile", "footer")
}

// アイコンアップロードハンドラ
func ProfileIconHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		http.Error(w, "セッションが無効です", http.StatusUnauthorized)
		return
	}

	// URLからユーザーIDを取得
	path := r.URL.Path
	targetUserID := path[len("/profile/icon/"):]

	// 自分のプロフィール以外での変更を防止
	if targetUserID != "" && targetUserID != session.User.ID {
		http.Error(w, "他のユーザーのアイコンは変更できません", http.StatusForbidden)
		return
	}

	// マルチパートフォームの解析
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// アイコンファイルを取得
	file, header, err := r.FormFile("icon")
	if err != nil {
		http.Error(w, "アイコンファイルの取得に失敗しました", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ファイルの拡張子を取得
	ext := filepath.Ext(header.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		http.Error(w, "サポートされていないファイル形式です", http.StatusBadRequest)
		return
	}

	// 一時ファイルを作成
	tempFile, err := os.CreateTemp("", "icon-*"+ext)
	if err != nil {
		http.Error(w, "一時ファイルの作成に失敗しました", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// ファイルをコピー
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "ファイルの保存に失敗しました", http.StatusInternalServerError)
		return
	}

	// 一時ファイルのパスを取得
	tempFilePath := tempFile.Name()

	// Firebase Storageにアップロード
	iconURL, err := firebase.UploadIcon(session.User.ID, tempFilePath)
	if err != nil {
		log.Printf("アイコンアップロードエラー: %v", err)
		http.Error(w, fmt.Sprintf("アイコンのアップロードに失敗しました: %v", err), http.StatusInternalServerError)
		return
	}

	// 一時ファイルを削除
	os.Remove(tempFilePath)

	// ユーザードキュメントを更新
	err = firebase.UpdateField("users", session.User.ID, "Icon", iconURL)
	if err != nil {
		http.Error(w, "ユーザー情報の更新に失敗しました", http.StatusInternalServerError)
		return
	}

	// プロフィールページにリダイレクト
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
