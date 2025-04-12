package repository

import (
	"context"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
)

// メールアドレスでユーザーを検索する
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

	// ユーザーが見つからない場合
	if len(docs) == 0 {
		return nil, nil
	}

	var user domain.User
	if err := docs[0].DataTo(&user); err != nil {
		return nil, err
	}

	// ドキュメントIDをユーザーIDとして設定
	user.ID = docs[0].Ref.ID

	return &user, nil
}

// ユーザーIDからユーザー情報を取得する
func GetUserByID(userID string) (*domain.User, error) {
	client, err := firebase.InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	doc, err := client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	user.ID = doc.Ref.ID

	return &user, nil
}
