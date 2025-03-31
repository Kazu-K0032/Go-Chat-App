package service

import (
	"fmt"
	"net/http"

	"security_chat_app/repository"
)

// Tweet ツイート情報を持つ構造体
type Tweet struct {
	Content string
	Date    string
}

// ProfileData プロフィールページのデータ構造体
type ProfileData struct {
	IsLoggedIn bool
	User       *repository.User
	Tweets     []Tweet
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Println("セッションの検証に成功しました")

	// ユーザー情報の取得
	user, err := repository.GetUserByEmail(session.Email)
	if err != nil {
		http.Error(w, "ユーザー情報の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// プロフィールデータの作成
	data := TemplateData{
		IsLoggedIn: true,
		User:       user,
		Tweets: []Tweet{
			{
				Content: "VBAの資格範囲終わった...。\n最初の章がデバッグに関する内容だったけど、最初にやるべきだったな。\nこれ知ってるだけで無駄な手間省けた。データ型を調べたりステップ実行で\nエラー特定したりと初期段階からやれてたらなあ",
				Date:    "2月24日",
			},
		},
	}

	// テンプレートを描画
	GenerateHTML(w, data, "layout", "header", "profile", "footer")
}
