package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"security_chat_app/config"
)

// layout.htmlをベースとしたHTMLを生成し、レスポンスに書きだす
func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

// HTTPサーバーを起動 + '/'へのアクセスでtopが実行されるようにする
func StartMainServer() error {
	http.HandleFunc("/", top)
	return http.ListenAndServe(":"+config.Config.Port, nil)
}
