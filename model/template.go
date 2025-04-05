package model

import (
	"security_chat_app/repository"
	"security_chat_app/service"
)

// TemplateData 共通のテンプレートデータ構造体
type TemplateData struct {
	IsLoggedIn       bool
	User             *repository.User
	SignupForm       SignupForm
	LoginForm        service.LoginForm
	ValidationErrors []string
}

// SignupForm 登録フォームのデータ構造体
type SignupForm struct {
	Name     string
	Email    string
	Password string
}
