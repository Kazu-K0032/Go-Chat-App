package handler

import (
	"fmt"
	"net/http"
	"time"

	"security_chat_app/internal/interface/presenter/html"
	"security_chat_app/repository"
)

type SignupForm struct {
	Name     string
	Email    string
	Password string
}

// SignupHandler 新規登録画面の表示と確認画面への遷移を処理
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := html.TemplateData{
			IsLoggedIn: false,
		}
		html.GenerateHTML(w, data, "layout", "header", "register", "footer")
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
			data := html.TemplateData{
				IsLoggedIn: false,
				SignupForm: form,
				Error:      "このメールアドレスは既に登録されています",
			}
			html.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		// デバッグ用出力
		fmt.Println("== 登録フォームの内容 ==")
		fmt.Println("名前:", form.Name)
		fmt.Println("メール:", form.Email)
		fmt.Println("パスワード:", form.Password)

		data := html.TemplateData{
			IsLoggedIn: false,
			SignupForm: form,
		}
		html.GenerateHTML(w, data, "layout", "header", "register_confirm", "footer")
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
		user := &repository.User{
			ID:        repository.GenerateUUID(),
			Name:      form.Name,
			Email:     form.Email,
			Password:  form.Password, // 実際の実装ではハッシュ化が必要
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Firestoreにユーザーを保存
		err := repository.AddData("users", user, user.ID)
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
