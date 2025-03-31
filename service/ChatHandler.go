package service

import (
	"fmt"
	"net/http"
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
	IsLoggedIn bool
	User       *repository.User
	ChatID     string
	TargetUser *repository.User
	Messages   []map[string]interface{}
}

// チャットページのハンドラ
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// URLからチャットIDを取得
	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "チャットIDが指定されていません", http.StatusBadRequest)
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

	// 対象ユーザーの情報を取得
	targetUser, err := GetUserData(targetUserID)
	if err != nil {
		http.Error(w, "対象ユーザーの情報の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// メッセージを取得
	messages, err := repository.GetChatMessages(chatID)
	if err != nil {
		http.Error(w, "メッセージの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// チャットページのデータを取得
	data := ChatPageData{
		IsLoggedIn: true,
		User:       session.User,
		ChatID:     chatID,
		TargetUser: targetUser,
		Messages:   messages,
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
	chatID := r.FormValue("chat_id")
	content := r.FormValue("content")

	if chatID == "" || content == "" {
		http.Error(w, "チャットIDとメッセージ内容が必要です", http.StatusBadRequest)
		return
	}

	// メッセージを作成
	message := map[string]interface{}{
		"user_id":    session.User.ID,
		"content":    content,
		"created_at": time.Now(),
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
	// チャットコレクションからデータを取得
	chatsData, err := repository.GetDataByQuery("chats", "participants", "array-contains", user.ID)
	if err != nil {
		return nil, err
	}

	var chats []Chat
	for _, data := range chatsData {
		// メッセージを取得
		messagesData, err := repository.GetDataByQuery("messages", "chatID", "==", data["id"])
		if err != nil {
			continue
		}

		var messages []Message
		for _, msgData := range messagesData {
			message := Message{
				ID:         msgData["id"].(string),
				Content:    msgData["content"].(string),
				SenderID:   msgData["senderID"].(string),
				SenderName: msgData["senderName"].(string),
				Time:       msgData["time"].(time.Time),
				IsRead:     msgData["isRead"].(bool),
			}
			messages = append(messages, message)
		}

		// 連絡先情報を取得
		contactData, err := repository.GetData("users", data["contactID"].(string))
		if err != nil {
			continue
		}

		contact := Contact{
			ID:       contactData["id"].(string),
			Username: contactData["username"].(string),
			Icon:     contactData["icon"].(string),
			LastSeen: contactData["lastSeen"].(time.Time),
		}

		chat := Chat{
			ID:        data["id"].(string),
			Contact:   contact,
			Messages:  messages,
			UpdatedAt: data["updatedAt"].(time.Time),
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

// メッセージIDを生成する
func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}
