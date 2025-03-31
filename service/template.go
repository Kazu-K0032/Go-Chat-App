package service

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"security_chat_app/repository"
)

// TemplateData テンプレートに渡すデータ構造体
type TemplateData struct {
	IsLoggedIn       bool
	SignupForm       SignupForm
	LoginForm        LoginForm
	Error            string
	Success          string
	ValidationErrors []string
	ResetForm        ResetForm
	User             *repository.User
	Tweets           []Tweet
}

// デフォルトアイコンのパス
const defaultIconPath = "/static/assets/default/user_icon"

// デフォルトアイコンの数
const defaultIconCount = 7

// デフォルトアイコンの名前
var defaultIconNames = []string{
	"elephant",
	"fox",
	"hamster",
	"koala",
	"monkey",
	"owl",
	"puma",
}

// ローカル乱数生成器
var localRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// テンプレートで使用する関数
var templateFuncs = template.FuncMap{
	"sub": func(a, b int) int {
		return a - b
	},
	"len": func(slice interface{}) int {
		switch v := slice.(type) {
		case []Message:
			return len(v)
		case []Contact:
			return len(v)
		case []Chat:
			return len(v)
		default:
			return 0
		}
	},
	"substr": func(s string, start, length int) string {
		if start < 0 {
			start = 0
		}
		if length < 0 {
			length = len(s)
		}
		if start > len(s) {
			return ""
		}
		end := start + length
		if end > len(s) {
			end = len(s)
		}
		return s[start:end]
	},
	"getRandomDefaultIcon": func() string {
		// 0から6までのランダムな数字を生成
		randomNum := localRand.Intn(defaultIconCount)
		// デフォルトアイコンのパスを生成
		return fmt.Sprintf("%s/default_icon_%s.png", defaultIconPath, defaultIconNames[randomNum])
	},
}

// GenerateHTML layout.htmlをベースとしたHTMLを生成し、レスポンスに書きだす
func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		path := fmt.Sprintf("app/templates/%s.html", file)
		files = append(files, path)
	}

	// テンプレートパース時のエラーハンドリング
	templates, err := template.New("layout").Funcs(templateFuncs).ParseFiles(files...)
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
