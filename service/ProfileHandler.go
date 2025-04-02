package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"security_chat_app/repository"
)

// Tweet ツイート情報を持つ構造体
type Tweet struct {
	Content string
	Date    string
}

// ProfileData プロフィールページのデータ構造体
type ProfileData struct {
	IsLoggedIn bool
	User       *repository.User
	Tweets     []Tweet
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

	// 最終更新日時を現在時刻に更新
	user.UpdatedAt = time.Now()
	err = repository.UpdateField("users", user.ID, "updated_at", user.UpdatedAt)
	if err != nil {
		http.Error(w, "最終更新日時の更新に失敗しました", http.StatusInternalServerError)
		return
	}

	// プロフィールデータの作成
	data := TemplateData{
		IsLoggedIn: true,
		User:       user,
		Tweets: []Tweet{
			{
				Content: "VBAの資格範囲終わった...。\n最初の章がデバッグに関する内容だったけど、最初にやるべきだったな。\nこれ知ってるだけで無駄な手間省けた。データ型を調べたりステップ実行で\nエラー特定したりと初期段階からやれてたらなあ",
				Date:    "2月24日",
			},
		},
		Replies: []Tweet{}, // 空のスライスで初期化
		Likes:   []Tweet{}, // 空のスライスで初期化
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
