package service

import (
	"fmt"
	"net/http"

	"security_chat_app/repository"
)

// 設定ページのデータ構造体
type SettingsPageData struct {
	IsLoggedIn bool
	User       *repository.User
}

// 設定ページのハンドラ
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 設定ページのデータを取得
	data, err := getSettingsPageData(session.User)
	if err != nil {
		http.Error(w, "設定ページのデータの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// テンプレートのレンダリング
	GenerateHTML(w, data, "layout", "header", "settings", "footer")
}

// 設定ページのデータを取得
func getSettingsPageData(user *repository.User) (SettingsPageData, error) {
	if user == nil {
		return SettingsPageData{}, fmt.Errorf("ユーザー情報が無効です")
	}

	return SettingsPageData{
		IsLoggedIn: true,
		User:       user,
	}, nil
}
