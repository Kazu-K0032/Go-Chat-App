package handler

import (
	"fmt"
	"net/http"

	"security_chat_app/internal/interface/html"
	"security_chat_app/internal/interface/middleware"
	"security_chat_app/internal/infrastructure/repository"
)

// LoginForm ログインフォームのデータ構造体
type LoginForm struct {
	Email    string
	Password string
}

// LoginHandler ログイン処理を実行
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := html.TemplateData{
			IsLoggedIn: false,
		}
		html.GenerateHTML(w, data, "layout", "header", "login", "footer")
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		form := LoginForm{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// バリデーション
		var validationErrors []string
		if form.Email == "" {
			validationErrors = append(validationErrors, "メールアドレスを入力してください")
		}
		if form.Password == "" {
			validationErrors = append(validationErrors, "パスワードを入力してください")
		}

		if len(validationErrors) > 0 {
			data := html.TemplateData{
				IsLoggedIn:       false,
				LoginForm:        html.LoginForm{Email: form.Email, Password: form.Password},
				ValidationErrors: validationErrors,
			}
			html.GenerateHTML(w, data, "layout", "header", "login", "footer")
			return
		}

		// ユーザー認証
		user, err := repository.GetUserByEmail(form.Email)
		if err != nil {
			data := html.TemplateData{
				IsLoggedIn:       false,
				LoginForm:        html.LoginForm{Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"認証エラーが発生しました"},
			}
			html.GenerateHTML(w, data, "layout", "header", "login", "footer")
			return
		}

		if user == nil || user.Password != form.Password {
			data := html.TemplateData{
				IsLoggedIn:       false,
				LoginForm:        html.LoginForm{Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"メールアドレスまたはパスワードが誤っています"},
			}
			html.GenerateHTML(w, data, "layout", "header", "login", "footer")
			return
		}

		// セッションの作成
		session, err := middleware.CreateSession(user)
		if err != nil {
			fmt.Println("セッション作成エラー:", err)
			data := html.TemplateData{
				IsLoggedIn:       false,
				LoginForm:        html.LoginForm{Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"セッション作成エラーが発生しました"},
			}
			html.GenerateHTML(w, data, "layout", "header", "login", "footer")
			return
		}

		fmt.Println("セッション作成成功:", session)

		// セッションクッキーの設定
		middleware.SetSessionCookie(w, session)
		fmt.Println("セッションクッキー設定完了")

		// ログイン成功後、プロフィールページにリダイレクト
		fmt.Println("リダイレクト開始: /profile")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
}
