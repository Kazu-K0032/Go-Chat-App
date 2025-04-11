package chat

import (
	"fmt"
	"sort"
	"time"

	"security_chat_app/internal/domain"
)

// チャットのコントローラー
type ChatController struct {
	chatUsecase domain.ChatUsecase
}

// chatUsecaseImpl はチャットのユースケースの実装
type chatUsecaseImpl struct {
	repo interface {
		AddChat(user, message string) error
	}
}

// NewChatUsecase はチャットのユースケースの実装を生成する
func NewChatUsecase(repo interface {
	AddChat(user, message string) error
},
) domain.ChatUsecase {
	return &chatUsecaseImpl{repo: repo}
}

// CreateChat はチャット作成時のビジネスロジックを定義
func (c *chatUsecaseImpl) CreateChat(user, message string) error {
	return c.repo.AddChat(user, message)
}

// 連絡先を交換したユーザーのデータを取得
func getContacts(user *repository.User) ([]domain.Contact, error) {
	// 連絡先コレクションからデータを取得
	contactsData, err := repository.GetDataByQuery("contacts", "userID", "==", user.ID)
	if err != nil {
		return nil, err
	}

	var contacts []domain.Contact
	for _, data := range contactsData {
		contact := domain.Contact{
			ID:       data["contactID"].(string),
			Username: data["username"].(string),
			IconURL:  data["iconURL"].(string),
			LastSeen: data["lastSeen"].(time.Time),
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// チャット履歴を取得
func GetChatHistory(user *repository.User) ([]domain.Chat, error) {
	// チャット履歴を取得
	chats, err := repository.GetAllData("chats")
	if err != nil {
		return nil, fmt.Errorf("チャット履歴の取得に失敗しました: %v", err)
	}

	var chatHistory []domain.Chat
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
		var messages []domain.Message
		var lastMessageTime time.Time
		for _, msg := range messagesData {
			messageTime := msg["created_at"].(time.Time)
			message := domain.Message{
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
		chatHistory = append(chatHistory, domain.Chat{
			ID: chatID,
			Contact: domain.Contact{
				ID:       targetUser.ID,
				Username: targetUser.Name,
				IconURL:  targetUser.IconURL,
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
