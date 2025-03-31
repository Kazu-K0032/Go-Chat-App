package router

import (
	"fmt"
	"net/http"

	"security_chat_app/repository"
	"security_chat_app/service"
)

// 共通のテンプレートデータ構造体
type TemplateData struct {
	IsLoggedIn bool
	User       *repository.User
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
	// セッションの検証
	session, err := service.ValidateSession(w, r)
	if err != nil {
		fmt.Println("セッションが無効です")
		data := TemplateData{
			IsLoggedIn: false,
		}
		service.GenerateHTML(w, data, "layout", "header", "top", "footer")
		return
	}

	// セッションが有効な場合
	data := TemplateData{
		IsLoggedIn: true,
		User:       session.User,
	}
	service.GenerateHTML(w, data, "layout", "header", "top", "footer")
}
