package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

type ChatRepository interface {
	AddChat(user, message string) error
}

type chatRepository struct {
	client *firestore.Client
}

func NewChatRepository(client *firestore.Client) ChatRepository {
	return &chatRepository{client: client}
}

// チャットの追加
func (r *chatRepository) AddChat(user, message string) error {
	_, _, err := r.client.Collection("chats").Add(context.Background(), map[string]interface{}{
		"user":    user,
		"message": message,
		"created": time.Now().Format("2025-03-28 00:00:00"),
	})
	return err
}
