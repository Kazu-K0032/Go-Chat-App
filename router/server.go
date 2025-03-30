package router

import (
	"fmt"
	"html/template"
	"net/http"

	"security_chat_app/config"
	"security_chat_app/service"
)

// layout.htmlをベースとしたHTMLを生成し、レスポンスに書きだす
func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

// HTTPサーバーを起動 + '/'へのアクセスでtopが実行されるようにする
func StartMainServer(chatUsecase service.ChatUsecase) error {
	// CSSファイル用
	cssFs := http.FileServer(http.Dir("app/css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFs))

	// JSファイル用
	jsFs := http.FileServer(http.Dir("app/js"))
	http.Handle("/js/", http.StripPrefix("/js/", jsFs))

	// チャット画面用
	http.HandleFunc("/", top)

	chatUsecase.CreateChat("test", "test")
	return http.ListenAndServe(":"+config.Config.Port, nil)
}
