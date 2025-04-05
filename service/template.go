package service

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"security_chat_app/repository"
)

// TemplateData 共通のテンプレートデータ構造体
type TemplateData struct {
	IsLoggedIn       bool
	User             *repository.User
	SignupForm       SignupForm
	LoginForm        LoginForm
	ResetForm        ResetForm
	ValidationErrors []string
	Error            string
}

// デフォルトアイコンのパス
const defaultIconPath = "icons/default"

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
		return fmt.Sprintf("%s/%s.png", defaultIconPath, defaultIconNames[randomNum])
	},
}

// GenerateHTML layout.htmlをベースとしたHTMLを生成し、レスポンスに書きだす
func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		path := fmt.Sprintf("app/templates/%s.html", file)
		files = append(files, path)
	}

	templates, err := template.New("layout").Funcs(templateFuncs).ParseFiles(files...)
	if err != nil {
		http.Error(writer, "テンプレートの読み込みに失敗しました", http.StatusInternalServerError)
		fmt.Println("テンプレート読み込みエラー:", err)
		return
	}

	// テンプレートをバッファに出力
	var buf bytes.Buffer
	err = templates.ExecuteTemplate(&buf, "layout", data)
	if err != nil {
		http.Error(writer, "テンプレートの実行に失敗しました", http.StatusInternalServerError)
		fmt.Println("テンプレート実行エラー:", err)
		return
	}

	// 成功したらまとめて出力
	buf.WriteTo(writer)
}
