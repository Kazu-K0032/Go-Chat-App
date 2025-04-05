package service

import (
	"fmt"
	"net/http"

	"security_chat_app/repository"
)

// LoginForm ログインフォームのデータ構造体
type LoginForm struct {
	Email    string
	Password string
}

// LoginHandler ログイン処理を実行
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := TemplateData{
			IsLoggedIn: false,
		}
		GenerateHTML(w, data, "layout", "header", "login", "footer")
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
			data := TemplateData{
				IsLoggedIn:       false,
				LoginForm:        form,
				ValidationErrors: validationErrors,
			}
			GenerateHTML(w, data, "layout", "header", "login", "footer")
			return
		}

		// ユーザー認証
		user, err := repository.GetUserByEmail(form.Email)
		if err != nil {
			data := TemplateData{
				IsLoggedIn:       false,
				LoginForm:        form,
				ValidationErrors: []string{"認証エラーが発生しました"},
			}
			GenerateHTML(w, data, "layout", "header", "login", "footer")
			return
		}

		if user == nil || user.Password != form.Password {
			data := TemplateData{
				IsLoggedIn:       false,
				LoginForm:        form,
				ValidationErrors: []string{"メールアドレスまたはパスワードが誤っています"},
			}
			GenerateHTML(w, data, "layout", "header", "login", "footer")
			return
		}

		// セッションの作成
		session, err := CreateSession(user)
		if err != nil {
			fmt.Println("セッション作成エラー:", err)
			data := TemplateData{
				IsLoggedIn:       false,
				LoginForm:        form,
				ValidationErrors: []string{"セッション作成エラーが発生しました"},
			}
			GenerateHTML(w, data, "layout", "header", "login", "footer")
			return
		}

		fmt.Println("セッション作成成功:", session)

		// セッションクッキーの設定
		SetSessionCookie(w, session)
		fmt.Println("セッションクッキー設定完了")

		// ログイン成功後、プロフィールページにリダイレクト
		fmt.Println("リダイレクト開始: /profile")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
}

// LogoutHandler ログアウト処理を実行
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := DeleteSession(w, r)
		if err != nil {
			http.Error(w, "ログアウトエラー", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
