package handler

import (
	"fmt"
	"net/http"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/interface/markup"
	"security_chat_app/internal/interface/middleware"
)

// 設定ページのデータ構造体
type SettingsPageData struct {
	IsLoggedIn       bool
	User             *domain.User
	ShowPasswordForm bool
	PasswordForm     struct {
		CurrentPassword    string
		NewPassword        string
		NewPasswordConfirm string
	}
	ValidationErrors []string
}

// 設定ページのハンドラ
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// パスワード変更フォームの処理
		r.ParseForm()
		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")
		newPasswordConfirm := r.FormValue("new_password_confirm")

		// バリデーション
		var validationErrors []string
		if currentPassword == "" {
			validationErrors = append(validationErrors, "現在のパスワードを入力してください")
		}
		if newPassword == "" {
			validationErrors = append(validationErrors, "新しいパスワードを入力してください")
		}
		if newPassword != newPasswordConfirm {
			validationErrors = append(validationErrors, "新しいパスワードが一致しません")
		}

		if len(validationErrors) > 0 {
			data := SettingsPageData{
				IsLoggedIn:       true,
				User:             session.User,
				ShowPasswordForm: true,
				PasswordForm: struct {
					CurrentPassword    string
					NewPassword        string
					NewPasswordConfirm string
				}{
					CurrentPassword:    currentPassword,
					NewPassword:        newPassword,
					NewPasswordConfirm: newPasswordConfirm,
				},
				ValidationErrors: validationErrors,
			}
			markup.GenerateHTML(w, data, "layout", "header", "settings", "footer")
			return
		}

		// パスワード更新
		err = firebase.UpdateField("users", session.User.ID, "password", newPassword)
		if err != nil {
			validationErrors = append(validationErrors, "パスワードの更新に失敗しました")
			data := SettingsPageData{
				IsLoggedIn:       true,
				User:             session.User,
				ShowPasswordForm: true,
				PasswordForm: struct {
					CurrentPassword    string
					NewPassword        string
					NewPasswordConfirm string
				}{
					CurrentPassword:    currentPassword,
					NewPassword:        newPassword,
					NewPasswordConfirm: newPasswordConfirm,
				},
				ValidationErrors: validationErrors,
			}
			markup.GenerateHTML(w, data, "layout", "header", "settings", "footer")
			return
		}

		// 成功時は設定ページにリダイレクト
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	// 設定ページのデータを取得
	data, err := getSettingsPageData(session.User, r)
	if err != nil {
		http.Error(w, "設定ページのデータの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// テンプレートのレンダリング
	markup.GenerateHTML(w, data, "layout", "header", "settings", "footer")
}

// 設定ページのデータを取得
func getSettingsPageData(user *domain.User, r *http.Request) (SettingsPageData, error) {
	if user == nil {
		return SettingsPageData{}, fmt.Errorf("ユーザー情報が無効です")
	}

	// クエリパラメータからフォームの表示状態を取得
	showPasswordForm := r.URL.Query().Get("show_password_form") == "true"

	return SettingsPageData{
		IsLoggedIn:       true,
		User:             user,
		ShowPasswordForm: showPasswordForm,
	}, nil
}
