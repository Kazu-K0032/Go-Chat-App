package user

import (
	"context"
	"time"

	"security_chat_app/repository"

	"golang.org/x/crypto/bcrypt"
)

// ユーザー登録
func CreateUser(name, email, password string) (*User, error) {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// ユーザーを作成
	user := &User{
		ID:        generateUUID(), // UUIDを生成する関数は別途実装が必要
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Firestoreにユーザーを保存
	client, err := repository.InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

