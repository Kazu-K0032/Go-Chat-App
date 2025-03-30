package router

import (
	"fmt"
	"net/http"

	"security_chat_app/service"
)

// 共通のテンプレートデータ構造体
type TemplateData struct {
	IsLoggedIn bool
	User       service.User
	SignupForm SignupForm
}

// 登録フォームのデータ構造体
type SignupForm struct {
	Name     string
	Email    string
	Password string
}

// '/'へのアクセス
func top(w http.ResponseWriter, r *http.Request) {
	_, err := service.ValidateSession(w, r)
	data := TemplateData{
		IsLoggedIn: err == nil,
	}
	if err != nil {
		fmt.Println("セッションが無効です")
		service.GenerateHTML(w, data, "layout", "header", "top", "footer")
	} else {
		fmt.Println("セッションが有効です")
		http.Redirect(w, r, "/chat", http.StatusFound)
	}
}

// ログインページ
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// ログイン処理
		// TODO: ログイン処理の実装
	}
	data := TemplateData{
		IsLoggedIn: false,
	}
	service.GenerateHTML(w, data, "layout", "header", "login", "footer")
}

// 登録処理
func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		form := SignupForm{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		data := TemplateData{
			IsLoggedIn: false,
			SignupForm: form,
		}
		service.GenerateHTML(w, data, "layout", "header", "register_confirm", "footer")
		return
	}

	data := TemplateData{
		IsLoggedIn: false,
	}
	service.GenerateHTML(w, data, "layout", "header", "register", "footer")
}

// 登録確認処理
func signupConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, err := service.CreateUser(name, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// プロフィールページ
func profile(w http.ResponseWriter, r *http.Request) {
	_, _ = service.ValidateSession(w, r)
	data := TemplateData{
		IsLoggedIn: true,
	}
	service.GenerateHTML(w, data, "layout", "header", "profile", "footer")
}

// 設定ページ
func settings(w http.ResponseWriter, r *http.Request) {
	_, _ = service.ValidateSession(w, r)
	data := TemplateData{
		IsLoggedIn: true,
	}
	service.GenerateHTML(w, data, "layout", "header", "settings", "footer")
}
