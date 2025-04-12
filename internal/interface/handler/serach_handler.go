package handler

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/interface/markup"
	"security_chat_app/internal/interface/middleware"
)

// 検索ページのデータ構造体
type SearchPageData struct {
	IsLoggedIn bool
	User       *domain.User
	Query      string
	Users      []map[string]interface{}
}

// 検索ハンドラ
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 検索ページのデータを取得
	data, err := getSearchPageData(session.User, r)
	if err != nil {
		http.Error(w, "検索データの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// テンプレートのレンダリング
	markup.GenerateHTML(w, data, "layout", "header", "search", "footer")
}

// 検索ページのデータを取得
func getSearchPageData(user *domain.User, r *http.Request) (SearchPageData, error) {
	if user == nil {
		return SearchPageData{}, fmt.Errorf("ユーザー情報が無効です")
	}

	// 検索クエリを取得
	query := r.URL.Query().Get("username")
	var users []map[string]interface{}
	var err error

	// 検索クエリがある場合は検索を実行、ない場合は全ユーザーを取得
	if query != "" {
		users, err = SearchUsers(query)
	} else {
		users, err = firebase.GetAllData("users")
	}

	if err != nil {
		return SearchPageData{}, fmt.Errorf("ユーザー情報の取得に失敗しました: %v", err)
	}

	// チャット履歴を取得
	chats, err := firebase.GetAllData("chats")
	if err != nil {
		return SearchPageData{}, fmt.Errorf("チャット履歴の取得に失敗しました: %v", err)
	}

	// チャット履歴のあるユーザーIDを集める
	chattedUsers := make(map[string]bool)
	for _, chatData := range chats {
		participants, ok := chatData["participants"].([]interface{})
		if !ok {
			continue
		}
		for _, p := range participants {
			if participantID, ok := p.(string); ok {
				chattedUsers[participantID] = true
			}
		}
	}

	// 自分以外かつチャット履歴のないユーザーをフィルタリング
	var filteredUsers []map[string]interface{}
	for _, u := range users {
		userID, ok := u["ID"].(string)
		if !ok {
			log.Printf("ユーザーIDの取得に失敗: %+v", u)
			continue
		}

		if userID != user.ID && !chattedUsers[userID] {
			// テンプレートで使用するフィールド名に合わせてデータを整形
			userData := map[string]interface{}{
				"id":       userID,
				"name":     u["Name"],
				"icon":     u["Icon"],
				"isOnline": u["IsOnline"],
			}
			log.Printf("フィルタリング後のユーザーデータ: %+v", userData)
			filteredUsers = append(filteredUsers, userData)
		}
	}

	// ユーザーをcreated_atで降順にソート
	sort.Slice(filteredUsers, func(i, j int) bool {
		timeI, okI := filteredUsers[i]["created_at"].(time.Time)
		timeJ, okJ := filteredUsers[j]["created_at"].(time.Time)

		// created_atがnilまたはtime.Timeでない場合は、最後に配置
		if !okI {
			return false
		}
		if !okJ {
			return true
		}

		return timeI.After(timeJ)
	})

	// 検索ページのデータを取得
	data := SearchPageData{
		IsLoggedIn: true,
		User:       user,
		Query:      query,
		Users:      filteredUsers,
	}

	return data, nil
}

// ユーザーを検索
func SearchUsers(query string) ([]map[string]interface{}, error) {
	// ユーザーを検索
	users, err := firebase.SearchUser(query)
	if err != nil {
		return nil, fmt.Errorf("ユーザーの検索に失敗しました: %v", err)
	}

	// 検索結果をログ出力
	log.Printf("検索結果: %+v", users)
	return users, nil
}

// ユーザー情報を取得
func GetUserData(userID string) (*domain.User, error) {
	// ユーザー情報を取得
	userData, err := firebase.GetData("users", userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザー情報の取得に失敗しました: %v", err)
	}

	// 必要なフィールドの存在確認と型アサーション
	id, ok := userData["ID"].(string)
	if !ok {
		return nil, fmt.Errorf("ユーザーIDの取得に失敗しました")
	}

	name, ok := userData["Name"].(string)
	if !ok {
		return nil, fmt.Errorf("ユーザー名の取得に失敗しました")
	}

	email, ok := userData["Email"].(string)
	if !ok {
		return nil, fmt.Errorf("メールアドレスの取得に失敗しました")
	}

	// アイコンURLを取得（存在しない場合は空文字列）
	iconURL := ""
	if icon, ok := userData["Icon"].(string); ok {
		iconURL = icon
	}

	// オンラインステータスを取得（存在しない場合はfalse）
	isOnline := false
	if online, ok := userData["IsOnline"].(bool); ok {
		isOnline = online
	}

	return &domain.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Icon:     iconURL,
		IsOnline: isOnline,
	}, nil
}
