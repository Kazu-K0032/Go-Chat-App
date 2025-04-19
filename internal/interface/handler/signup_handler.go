package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/interface/markup"
	"security_chat_app/internal/interface/middleware"
	utils "security_chat_app/internal/utils/uuid"
)

// SignupHandler 新規登録画面の表示と確認画面への遷移を処理
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// サインアップ画面の表示
	if r.Method == http.MethodGet {
		data := domain.TemplateData{
			IsLoggedIn: false,
		}
		markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
		return
	}

	// サインアップ処理
	if r.Method == http.MethodPost {
		r.ParseForm()
		form := domain.SignupForm{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// バリデーション
		var validationErrors []string
		if form.Name == "" {
			validationErrors = append(validationErrors, "名前を入力してください")
		}
		if form.Email == "" {
			validationErrors = append(validationErrors, "メールアドレスを入力してください")
		}
		if !strings.Contains(form.Email, "@") {
			validationErrors = append(validationErrors, "有効なメールアドレスを入力してください")
		}
		if form.Password == "" {
			validationErrors = append(validationErrors, "パスワードを入力してください")
		}
		if len(form.Password) < 8 {
			validationErrors = append(validationErrors, "パスワードは8文字以上で入力してください")
		}

		if len(validationErrors) > 0 {
			log.Printf("バリデーションエラー: %v", validationErrors)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       domain.SignupForm{Name: form.Name, Email: form.Email, Password: form.Password},
				ValidationErrors: validationErrors,
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		// メールアドレスの重複チェック
		existingUsers, err := firebase.GetDataByQuery("users", "Email", "==", form.Email)
		if err != nil {
			log.Printf("ユーザー検索エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       domain.SignupForm{Name: form.Name, Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		if len(existingUsers) > 0 {
			log.Printf("メールアドレス重複エラー: %s", form.Email)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       domain.SignupForm{Name: form.Name, Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"このメールアドレスは既に登録されています"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		// パスワードのハッシュ化
		hashedPassword, err := utils.HashPassword(form.Password)
		if err != nil {
			log.Printf("パスワードハッシュ化エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       domain.SignupForm{Name: form.Name, Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		// UUIDの生成
		userID, err := utils.GenerateUUID()
		if err != nil {
			log.Printf("UUID生成エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       domain.SignupForm{Name: form.Name, Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		// ユーザーの作成
		user := &domain.User{
			ID:        userID,
			Name:      form.Name,
			Email:     form.Email,
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Firestoreにユーザーを保存
		err = firebase.AddData("users", user, user.ID)
		if err != nil {
			log.Printf("ユーザー作成エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       domain.SignupForm{Name: form.Name, Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"ユーザー作成エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		// セッションの作成
		session, err := middleware.CreateSession(user)
		if err != nil {
			log.Printf("セッション作成エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       domain.SignupForm{Name: form.Name, Email: form.Email, Password: form.Password},
				ValidationErrors: []string{"セッション作成エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		// セッションクッキーの設定
		middleware.SetSessionCookie(w, session)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// その他のHTTPメソッドは許可しない
	log.Printf("不正なHTTPメソッド: %s", r.Method)
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// SignupConfirmHandler 登録内容の確認とFirebaseへの保存を処理
func SignupConfirmHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		// GETメソッドの場合はフォームの値をURLから取得
		var form domain.SignupForm
		if r.Method == http.MethodGet {
			form = domain.SignupForm{
				Name:     r.URL.Query().Get("name"),
				Email:    r.URL.Query().Get("email"),
				Password: r.URL.Query().Get("password"),
			}
		} else {
			r.ParseForm()
			form = domain.SignupForm{
				Name:     r.FormValue("name"),
				Email:    r.FormValue("email"),
				Password: r.FormValue("password"),
			}
		}

		// バリデーション
		var validationErrors []string
		if form.Name == "" {
			validationErrors = append(validationErrors, "名前を入力してください")
		}
		if form.Email == "" {
			validationErrors = append(validationErrors, "メールアドレスを入力してください")
		}
		if !strings.Contains(form.Email, "@") {
			validationErrors = append(validationErrors, "有効なメールアドレスを入力してください")
		}
		if form.Password == "" {
			validationErrors = append(validationErrors, "パスワードを入力してください")
		}
		if len(form.Password) < 8 {
			validationErrors = append(validationErrors, "パスワードは8文字以上で入力してください")
		}

		if len(validationErrors) > 0 {
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       form,
				ValidationErrors: validationErrors,
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		// メールアドレスの重複チェック
		existingUsers, err := firebase.GetDataByQuery("users", "Email", "==", form.Email)
		if err != nil {
			log.Printf("ユーザー検索エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       form,
				ValidationErrors: []string{"エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		if len(existingUsers) > 0 {
			data := domain.TemplateData{
				IsLoggedIn:       false,
				SignupForm:       form,
				ValidationErrors: []string{"このメールアドレスは既に登録されています"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
			return
		}

		if r.Method == http.MethodPost {
			// UUIDの生成
			userID, err := utils.GenerateUUID()
			if err != nil {
				log.Printf("UUID生成エラー: %v", err)
				data := domain.TemplateData{
					IsLoggedIn:       false,
					SignupForm:       form,
					ValidationErrors: []string{"エラーが発生しました"},
				}
				markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
				return
			}

			// パスワードのハッシュ化
			hashedPassword, err := utils.HashPassword(form.Password)
			if err != nil {
				log.Printf("パスワードハッシュ化エラー: %v", err)
				data := domain.TemplateData{
					IsLoggedIn:       false,
					SignupForm:       form,
					ValidationErrors: []string{"エラーが発生しました"},
				}
				markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
				return
			}

			// ユーザーデータの作成
			user := &domain.User{
				ID:        userID,
				Name:      form.Name,
				Email:     form.Email,
				Password:  hashedPassword,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// Firestoreにユーザーを保存
			err = firebase.AddData("users", user, user.ID)
			if err != nil {
				log.Printf("ユーザー作成エラー: %v", err)
				data := domain.TemplateData{
					IsLoggedIn:       false,
					SignupForm:       form,
					ValidationErrors: []string{"ユーザー作成エラーが発生しました"},
				}
				markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
				return
			}

			// 登録成功後、ログインページにリダイレクト
			http.Redirect(w, r, "/login?success=true", http.StatusSeeOther)
			return
		}

		// 確認画面の表示
		data := domain.TemplateData{
			IsLoggedIn: false,
			SignupForm: form,
		}
		markup.GenerateHTML(w, data, "layout", "header", "register_confirm", "footer")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
