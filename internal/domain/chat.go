package domain

import (
	"time"

	"security_chat_app/repository"
)

// 連絡先を交換したユーザーの構造体
type Contact struct {
	ID       string
	Username string
	IconURL  string
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
