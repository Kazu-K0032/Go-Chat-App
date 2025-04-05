package router

import (
	"net/http"

	"security_chat_app/repository"
	"security_chat_app/service"
)

// 共通のテンプレートデータ構造体
type TemplateData struct {
	IsLoggedIn       bool
	User             *repository.User
	SignupForm       SignupForm
	LoginForm        service.LoginForm
	ValidationErrors []string
}

// 登録フォームのデータ構造体
type SignupForm struct {
	Name     string
	Email    string
	Password string
}

// '/'へのアクセス
func top(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	_, err := service.ValidateSession(w, r)
	if err != nil {
		// 未ログインの場合はログイン画面を表示
		data := service.TemplateData{
			IsLoggedIn: false,
		}
		service.GenerateHTML(w, data, "layout", "header", "login", "footer")
		return
	}

	// ログイン済みの場合は検索ページにリダイレクト
	http.Redirect(w, r, "/search", http.StatusSeeOther)
}
