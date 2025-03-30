package router

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"security_chat_app/service"
)

// contextKey コンテキストのキーとして使用するカスタム型
type contextKey string

const templateDataKey contextKey = "templateData"

// layout.htmlをベースとしたHTMLを生成し、レスポンスに書きだす
func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		path := fmt.Sprintf("app/templates/%s.html", file)
		files = append(files, path)
	}
	templates, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(writer, "テンプレートの読み込みに失敗しました", http.StatusInternalServerError)
		fmt.Println("テンプレート読み込みエラー:", err)
		return
	}

	err = templates.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		http.Error(writer, "テンプレートの実行に失敗しました", http.StatusInternalServerError)
		fmt.Println("テンプレート実行エラー:", err)
	}
}

// Middleware セッション管理のミドルウェア
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// セッションの検証
		_, err := service.ValidateSession(w, r)
		if err != nil {
			// セッションが無効な場合は、ログインしていない状態として処理
			data := service.TemplateData{IsLoggedIn: false}
			r = r.WithContext(context.WithValue(r.Context(), templateDataKey, data))
		} else {
			// セッションが有効な場合は、ログインしている状態として処理
			data := service.TemplateData{IsLoggedIn: true}
			r = r.WithContext(context.WithValue(r.Context(), templateDataKey, data))
		}

		next.ServeHTTP(w, r)
	})
}

// SetupRouter ルーティングの設定
func SetupRouter(chatUsecase service.ChatUsecase) *http.ServeMux {
	mux := http.NewServeMux()

	// 静的ファイル (CSS/JS)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("app/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("app/js"))))

	// 静的ファイルの提供
	fs := http.FileServer(http.Dir("app/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// ルーティングの設定
	mux.HandleFunc("/", top)
	mux.HandleFunc("/login", service.LoginHandler)
	mux.HandleFunc("/logout", service.LogoutHandler)
	mux.HandleFunc("/signup", service.SignupHandler)
	mux.HandleFunc("/signup/confirm", service.SignupConfirmHandler)

	return mux
}

func StartMainServer(chatUsecase service.ChatUsecase) error {
	mux := SetupRouter(chatUsecase)
	return http.ListenAndServe(":8080", mux) // <- ここで mux を渡す！
}
