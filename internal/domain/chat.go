package domain

import (
	"time"

	"security_chat_app/internal/usecase/chat"
)

// // 連絡先を交換したユーザーの構造体
// type Contact struct {
// 	ID       string
// 	Username string
// 	IconURL  string
// 	LastSeen time.Time
// }

// // 対象ユーザーとのチャット履歴を管理する構造体
// type Chat struct {
// 	ID        string
// 	Contact   Contact
// 	Messages  []Message
// 	UpdatedAt time.Time
// }

// // メッセージの構造体
// type Message struct {
// 	ID         string
// 	Content    string
// 	SenderID   string
// 	SenderName string
// 	Time       time.Time
// 	IsRead     bool
// }

// // チャットページのデータ構造体
// type ChatPageData struct {
// 	IsLoggedIn  bool
// 	User        *User
// 	ChatID      string
// 	TargetUser  *User
// 	Messages    []Message
// 	Chats       []Chat
// 	CurrentChat *Chat
// }

// チャットの構造体
type Chats struct {
	ID        string
	IsGroup   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// チャット参加者の構造体
type ChatParticipants struct {
	ID       string
	ChatID   string
	UserID   string
	JoinedAt time.Time
}

// メッセージの構造体
type Messages struct {
	ID         string
	ChatID     string
	SenderID   string
	SenderName string
	Content    string
	CreatedAt  time.Time
	IsRead     bool
}

// チャット開始ハンドラ
type ChatController struct {
	chatUsecase ChatUsecase
}

// チャットのユースケース
type ChatUsecase interface {
	CreateChat(user, message string) error
}

// チャットのユースケースの実装
type chatUsecase struct {
	repo interface {
		AddChat(user, message string) error
	}
}

// チャットのユースケースの実装のコンストラクタ
func NewChatUsecase(repo interface {
	AddChat(user, message string) error
},
) ChatUsecase {
	return &chat.chatUsecase{repo: repo}
}
