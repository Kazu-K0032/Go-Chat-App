package repository

import (
	"context"
	"log"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
)

// メールアドレスでユーザーを検索する
func GetUserByEmail(email string) (*domain.User, error) {
	client, err := firebase.InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := client.Collection("users").Where("email", "==", email)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("ユーザー検索エラー: %v, email=%s", err, email)
		return nil, err
	}

	// ユーザーが見つからない場合
	if len(docs) == 0 {
		log.Printf("ユーザーが見つかりません: email=%s", email)
		return nil, nil
	}

	var user domain.User
	if err := docs[0].DataTo(&user); err != nil {
		log.Printf("ユーザーデータ変換エラー: %v", err)
		return nil, err
	}

	// ドキュメントIDをユーザーIDとして設定
	user.ID = docs[0].Ref.ID
	log.Printf("ユーザーを取得しました: id=%s, email=%s", user.ID, user.Email)
	return &user, nil
}

// ユーザーIDからユーザー情報を取得する
func GetUserByID(userID string) (*domain.User, error) {
	client, err := firebase.InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doc, err := client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		log.Printf("ユーザー取得エラー: %v, userID=%s", err, userID)
		return nil, err
	}

	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		log.Printf("ユーザーデータ変換エラー: %v", err)
		return nil, err
	}
	user.ID = doc.Ref.ID

	log.Printf("ユーザーを取得しました: id=%s, email=%s", user.ID, user.Email)
	return &user, nil
}
