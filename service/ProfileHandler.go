// プロフィールページで使用する関数などを定義するファイルです

package service

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"security_chat_app/repository"
)

// ProfileData プロフィールページのデータ構造体
type ProfileData struct {
	IsLoggedIn bool
	User       *repository.User
	Posts      []repository.Post
	Replies    []repository.Post
	Likes      []repository.Post
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Println("セッションの検証に成功しました")

	// ユーザー情報の取得
	user, err := repository.GetUserByEmail(session.Email)
	if err != nil {
		http.Error(w, "ユーザー情報の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// アイコンが設定されていない場合はデフォルトアイコンを設定
	if user.IconURL == "" {
		// デフォルトアイコンのパスを生成
		randomNum := rand.Intn(7)
		defaultIconNames := []string{"elephant", "fox", "hamster", "koala", "monkey", "owl", "puma"}
		defaultIconPath := fmt.Sprintf("icons/default/default_icon_%s.png", defaultIconNames[randomNum])

		fmt.Printf("デフォルトアイコンを設定中: %s\n", defaultIconPath)

		// デフォルトアイコンのURLを取得
		iconURL, err := repository.GetDefaultIconURL(defaultIconPath)
		if err != nil {
			fmt.Printf("デフォルトアイコンの取得に失敗: %v\n", err)
			http.Error(w, "デフォルトアイコンの取得に失敗しました", http.StatusInternalServerError)
			return
		}

		// ユーザーのIconURLを更新
		user.IconURL = iconURL
		err = repository.UpdateField("users", user.ID, "iconURL", iconURL)
		if err != nil {
			http.Error(w, "アイコンURLの更新に失敗しました", http.StatusInternalServerError)
			return
		}

		fmt.Printf("デフォルトアイコンを設定しました: %s\n", iconURL)
	}

	// 投稿、返信、いいねを取得
	posts, err := repository.GetUserPosts(user.ID)
	if err != nil {
		fmt.Printf("投稿の取得に失敗: %v\n", err)
		posts = []repository.Post{}
	}

	replies, err := repository.GetUserReplies(user.ID)
	if err != nil {
		fmt.Printf("返信の取得に失敗: %v\n", err)
		replies = []repository.Post{}
	}

	likes, err := repository.GetUserLikedPosts(user.ID)
	if err != nil {
		fmt.Printf("いいねの取得に失敗: %v\n", err)
		likes = []repository.Post{}
	}

	// 最終更新日時を現在時刻に更新
	user.UpdatedAt = time.Now()
	err = repository.UpdateField("users", user.ID, "updated_at", user.UpdatedAt)
	if err != nil {
		http.Error(w, "最終更新日時の更新に失敗しました", http.StatusInternalServerError)
		return
	}

	// プロフィールデータの作成
	data := ProfileData{
		IsLoggedIn: true,
		User:       user,
		Posts:      posts,
		Replies:    replies,
		Likes:      likes,
	}

	// テンプレートを描画
	GenerateHTML(w, data, "layout", "header", "profile", "footer")
}

// アイコンアップロードハンドラ
func ProfileIconHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := ValidateSession(w, r)
	if err != nil {
		http.Error(w, "セッションが無効です", http.StatusUnauthorized)
		return
	}

	// マルチパートフォームの解析
	err = r.ParseMultipartForm(10 << 20) // 10MBの制限
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
	iconURL, err := repository.UploadIcon(session.User.ID, tempFilePath)
	if err != nil {
		http.Error(w, "アイコンのアップロードに失敗しました", http.StatusInternalServerError)
		return
	}

	// 一時ファイルを削除
	os.Remove(tempFilePath)

	// ユーザードキュメントを更新
	err = repository.UpdateField("users", session.User.ID, "iconURL", iconURL)
	if err != nil {
		http.Error(w, "ユーザー情報の更新に失敗しました", http.StatusInternalServerError)
		return
	}

	// プロフィールページにリダイレクト
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
