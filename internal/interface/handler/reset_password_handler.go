package handler

import (
	"fmt"
	"net/http"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/infrastructure/repository"
	"security_chat_app/internal/interface/markup"
)

// ResetForm パスワード再設定フォームのデータ構造体
type ResetForm struct {
	Email           string
	Password        string
	PasswordConfirm string
}

// ResetPasswordHandler パスワード再設定処理を実行
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := domain.TemplateData{
			IsLoggedIn: false,
			ResetForm:  domain.ResetForm{},
		}
		markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		form := ResetForm{
			Email:           r.FormValue("email"),
			Password:        r.FormValue("password"),
			PasswordConfirm: r.FormValue("password_confirm"),
		}

		// バリデーション
		var validationErrors []string
		if form.Email == "" {
			validationErrors = append(validationErrors, "メールアドレスを入力してください")
		}
		if form.Password == "" {
			validationErrors = append(validationErrors, "新しいパスワードを入力してください")
		}
		if form.Password != form.PasswordConfirm {
			validationErrors = append(validationErrors, "パスワードが一致しません")
		}

		if len(validationErrors) > 0 {
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        domain.ResetForm{Email: form.Email},
				ValidationErrors: validationErrors,
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		// ユーザー検索
		user, err := repository.GetUserByEmail(form.Email)
		if err != nil {
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        domain.ResetForm{Email: form.Email},
				ValidationErrors: []string{"ユーザー検索エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		if user == nil {
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        domain.ResetForm{Email: form.Email},
				ValidationErrors: []string{"該当するユーザーが見つかりません"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		// パスワード更新
		fmt.Println(user.ID)
		fmt.Println(form.Password)
		err = firebase.UpdateField("users", user.ID, "password", form.Password)
		if err != nil {
			fmt.Println("Firestore Update エラー:", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        domain.ResetForm{Email: form.Email},
				ValidationErrors: []string{"パスワード更新エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		fmt.Println("Firestore Update 成功")

		// 成功時はホームページにリダイレクト
		http.Redirect(w, r, "/?success=パスワードを再設定しました", http.StatusSeeOther)
	}
}
