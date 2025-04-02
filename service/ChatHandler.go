package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"security_chat_app/repository"
)

// 連絡先を交換したユーザーの構造体
type Contact struct {
	ID       string
	Username string
	Icon     string
	LastSeen time.Time
}

// 対象ユーザーとのチャット履歴を管理する構造体
type Chat struct {
	ID        string
	Contact   Contact
	Messages  []Message
	UpdatedAt time.Time
}

// メッセージの構造体
type Message struct {
	ID         string
	Content    string
	SenderID   string
	SenderName string
	Time       time.Time
	IsRead     bool
}

// チャットページのデータ構造体
type ChatPageData struct {
	IsLoggedIn  bool
	User        *repository.User
	ChatID      string
	TargetUser  *repository.User
	Messages    []Message
	Chats       []Chat
	CurrentChat *Chat
}

// チャットページのハンドラ
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// POSTリクエストの場合はメッセージ送信処理
	if r.Method == http.MethodPost {
		// フォームデータから情報を取得
		chatID := r.FormValue("chatID")
		content := r.FormValue("content")

		if chatID == "" || content == "" {
			http.Error(w, "チャットIDとメッセージ内容が必要です", http.StatusBadRequest)
			return
		}

		// メッセージを作成
		message := map[string]interface{}{
			"sender_id":   session.User.ID,
			"sender_name": session.User.Name,
			"content":     content,
			"created_at":  time.Now(),
			"is_read":     false,
		}

		// メッセージを保存
		err = repository.AddChatMessage(chatID, message)
		if err != nil {
			http.Error(w, "メッセージの送信に失敗しました", http.StatusInternalServerError)
			return
		}

		// JSONレスポンスを返す
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"content": content,
			"time":    time.Now().Format("15:04"),
		})
		return
	}

	// チャット一覧を取得
	chats, err := getChatHistory(session.User)
	if err != nil {
		http.Error(w, "チャット一覧の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// URLからチャットIDを取得
	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		// チャットIDがない場合は、チャット一覧を表示
		data := ChatPageData{
			IsLoggedIn: true,
			User:       session.User,
			Chats:      chats,
		}
		GenerateHTML(w, data, "layout", "header", "chat", "footer")
		return
	}

	// チャットの存在確認
	exists, err := repository.CheckChatExists(chatID)
	if err != nil {
		http.Error(w, "チャットの確認に失敗しました", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "チャットが見つかりません", http.StatusNotFound)
		return
	}

	// チャットの参加者を取得
	participants, err := repository.GetChatParticipants(chatID)
	if err != nil {
		http.Error(w, "チャットの参加者情報の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// 対象ユーザーを特定
	var targetUserID string
	for _, p := range participants {
		if p != session.User.ID {
			targetUserID = p
			break
		}
	}

	// 対象ユーザーの存在確認
	_, err = GetUserData(targetUserID)
	if err != nil {
		http.Error(w, "対象ユーザーが見つかりません", http.StatusNotFound)
		return
	}

	// 対象ユーザーの情報を取得
	targetUser, err := GetUserData(targetUserID)
	if err != nil {
		http.Error(w, "対象ユーザーの情報の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// メッセージを取得
	messagesData, err := repository.GetChatMessages(chatID)
	if err != nil {
		http.Error(w, "メッセージの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// メッセージの型変換
	var messages []Message
	for _, msg := range messagesData {
		message := Message{
			ID:         msg["id"].(string),
			Content:    msg["content"].(string),
			SenderID:   msg["sender_id"].(string),
			SenderName: msg["sender_name"].(string),
			Time:       msg["created_at"].(time.Time),
			IsRead:     msg["is_read"].(bool),
		}
		messages = append(messages, message)
	}

	// 現在のチャットを特定
	var currentChat *Chat
	for _, chat := range chats {
		if chat.ID == chatID {
			currentChat = &chat
			break
		}
	}

	// チャットページのデータを取得
	data := ChatPageData{
		IsLoggedIn:  true,
		User:        session.User,
		ChatID:      chatID,
		TargetUser:  targetUser,
		Messages:    messages,
		Chats:       chats,
		CurrentChat: currentChat,
	}

	// テンプレートのレンダリング
	GenerateHTML(w, data, "layout", "header", "chat", "footer")
}

// メッセージ送信ハンドラ
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// セッションの検証
	session, err := ValidateSession(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// フォームデータから情報を取得
	chatID := r.FormValue("chatID")
	content := r.FormValue("content")

	if chatID == "" || content == "" {
		http.Error(w, "チャットIDとメッセージ内容が必要です", http.StatusBadRequest)
		return
	}

	// メッセージを作成
	message := map[string]interface{}{
		"sender_id":   session.User.ID,
		"sender_name": session.User.Name,
		"content":     content,
		"created_at":  time.Now(),
		"is_read":     false,
	}

	// メッセージを保存
	err = repository.AddChatMessage(chatID, message)
	if err != nil {
		http.Error(w, "メッセージの送信に失敗しました", http.StatusInternalServerError)
		return
	}

	// チャットページにリダイレクト
	http.Redirect(w, r, fmt.Sprintf("/chat?chat_id=%s", chatID), http.StatusSeeOther)
}

// チャット開始ハンドラ
func StartChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// セッションの検証
	session, err := ValidateSession(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// URLから対象ユーザーIDを取得
	targetUserID := r.URL.Path[len("/chat/"):]
	if targetUserID == "" {
		http.Error(w, "ユーザーIDが指定されていません", http.StatusBadRequest)
		return
	}

	// 対象ユーザーの存在確認
	_, err = GetUserData(targetUserID)
	if err != nil {
		http.Error(w, "対象ユーザーが見つかりません", http.StatusNotFound)
		return
	}

	// チャットを開始
	chatID, err := repository.StartChat(session.User.ID, targetUserID)
	if err != nil {
		log.Printf("チャット開始エラー: %v", err)
		http.Error(w, "チャットの開始に失敗しました", http.StatusInternalServerError)
		return
	}

	log.Printf("チャットを開始しました: chatID=%s, userID=%s, targetUserID=%s",
		chatID, session.User.ID, targetUserID)

	// チャットページにリダイレクト
	redirectURL := fmt.Sprintf("/chat?chat_id=%s", chatID)
	log.Printf("リダイレクト先: %s", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// 連絡先を交換したユーザーのデータを取得
func getContacts(user *repository.User) ([]Contact, error) {
	// 連絡先コレクションからデータを取得
	contactsData, err := repository.GetDataByQuery("contacts", "userID", "==", user.ID)
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	for _, data := range contactsData {
		contact := Contact{
			ID:       data["contactID"].(string),
			Username: data["username"].(string),
			Icon:     data["icon"].(string),
			LastSeen: data["lastSeen"].(time.Time),
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// チャット履歴を取得
func getChatHistory(user *repository.User) ([]Chat, error) {
	// チャット履歴を取得
	chats, err := repository.GetAllData("chats")
	if err != nil {
		return nil, fmt.Errorf("チャット履歴の取得に失敗しました: %v", err)
	}

	var chatHistory []Chat
	seenChats := make(map[string]bool) // 重複チェック用のマップ

	for _, chatData := range chats {
		// チャットIDの取得
		chatID, ok := chatData["id"].(string)
		if !ok {
			continue
		}

		// 参加者の取得
		participants, ok := chatData["participants"].([]interface{})
		if !ok {
			continue
		}

		// 現在のユーザーが参加者かチェック
		isParticipant := false
		var targetUserID string
		for _, p := range participants {
			if participantID, ok := p.(string); ok {
				if participantID == user.ID {
					isParticipant = true
				} else {
					targetUserID = participantID
				}
			}
		}

		if !isParticipant || seenChats[chatID] {
			continue
		}
		seenChats[chatID] = true

		// メッセージの取得
		messagesData, err := repository.GetChatMessages(chatID)
		if err != nil {
			continue
		}

		// メッセージの型変換
		var messages []Message
		var lastMessageTime time.Time
		for _, msg := range messagesData {
			messageTime := msg["created_at"].(time.Time)
			message := Message{
				ID:         msg["id"].(string),
				Content:    msg["content"].(string),
				SenderID:   msg["sender_id"].(string),
				SenderName: msg["sender_name"].(string),
				Time:       messageTime,
				IsRead:     msg["is_read"].(bool),
			}
			messages = append(messages, message)

			// 最新のメッセージ時刻を更新
			if messageTime.After(lastMessageTime) {
				lastMessageTime = messageTime
			}
		}

		// チャット相手の情報を取得
		targetUser, err := GetUserData(targetUserID)
		if err != nil {
			continue
		}

		// チャット履歴に追加
		chatHistory = append(chatHistory, Chat{
			ID: chatID,
			Contact: Contact{
				ID:       targetUser.ID,
				Username: targetUser.Name,
				Icon:     targetUser.Icon,
				LastSeen: time.Now(), // 仮の値として現在時刻を設定
			},
			Messages:  messages,
			UpdatedAt: lastMessageTime,
		})
	}

	// 更新時刻でソート（新しい順）
	sort.Slice(chatHistory, func(i, j int) bool {
		return chatHistory[i].UpdatedAt.After(chatHistory[j].UpdatedAt)
	})

	return chatHistory, nil
}

// メッセージIDを生成する
func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}
