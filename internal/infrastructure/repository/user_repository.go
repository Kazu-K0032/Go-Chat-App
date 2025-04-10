package repository

import (
	"context"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
)

// GetUserByEmail メールアドレスでユーザーを検索する
func GetUserByEmail(email string) (*domain.User, error) {
	client, err := firebase.InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	query := client.Collection("users").Where("email", "==", email)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil // ユーザーが見つからない
	}

	var user domain.User
	if err := docs[0].DataTo(&user); err != nil {
		return nil, err
	}

	// ドキュメントIDをユーザーIDとして設定
	user.ID = docs[0].Ref.ID

	return &user, nil
}
