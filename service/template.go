package service

import (
	"fmt"
	"html/template"
	"net/http"
)

// TemplateData テンプレートに渡すデータ構造体
type TemplateData struct {
	IsLoggedIn       bool
	SignupForm       SignupForm
	LoginForm        LoginForm
	Error            string
	Success          string
	ValidationErrors []string
}

// GenerateHTML layout.htmlをベースとしたHTMLを生成し、レスポンスに書きだす
func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		path := fmt.Sprintf("app/templates/%s.html", file)
		files = append(files, path)
	}

	// テンプレートパース時のエラーハンドリング
	templates, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(writer, "テンプレートの読み込みに失敗しました", http.StatusInternalServerError)
		fmt.Println("テンプレート読み込みエラー:", err)
		return
	}

	// テンプレート実行時のエラーハンドリング
	err = templates.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		http.Error(writer, "テンプレートの実行に失敗しました", http.StatusInternalServerError)
		fmt.Println("テンプレート実行エラー:", err)
		fmt.Printf("渡されたデータ: %#v\n", data)
	}
}
