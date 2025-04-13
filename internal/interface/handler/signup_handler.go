package handler

import (
	"net/http"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/infrastructure/repository"
	"security_chat_app/internal/interface/markup"
	utils "security_chat_app/internal/utils/uuid"
)

type SignupForm struct {
	Name     string
	Email    string
	Password string
}

// SignupHandler 新規登録画面の表示と確認画面への遷移を処理
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := domain.TemplateData{
			IsLoggedIn: false,
		}
		markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		form := SignupForm{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// メールアドレスの重複チェック
		existingUser, err := repository.GetUserByEmail(form.Email)
		if err != nil {
			http.Error(w, "データベースエラー", http.StatusInternalServerError)
			return
		}
		if existingUser != nil {
			data := domain.TemplateData{
				IsLoggedIn: false,
				SignupForm: domain.SignupForm{
					Username: form.Name,
					Email:    form.Email,
					Password: form.Password,
				},
				Error: "このメールアドレスは既に登録されています",
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}
		data := domain.TemplateData{
			IsLoggedIn: false,
			SignupForm: domain.SignupForm{
				Username: form.Name,
				Email:    form.Email,
				Password: form.Password,
			},
		}
		markup.GenerateHTML(w, data, "layout", "header", "register_confirm", "footer")
	}
}

// SignupConfirmHandler 登録内容の確認とFirebaseへの保存を処理
func SignupConfirmHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		form := SignupForm{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// ユーザーデータの作成
		user := &domain.User{
			ID:        utils.GenerateUUID(),
			Name:      form.Name,
			Email:     form.Email,
			Password:  form.Password,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Firestoreにユーザーを保存
		err := firebase.AddData("users", user, user.ID)
		if err != nil {
			http.Error(w, "ユーザー登録エラー", http.StatusInternalServerError)
			return
		}

		// 登録成功後、ログインページにリダイレクト
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
