package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// User ユーザーモデル
type User struct {
	ID        string    `json:"id" firestore:"id"`
	Name      string    `json:"name" firestore:"name"`
	Email     string    `json:"email" firestore:"email"`
	Password  string    `json:"password" firestore:"password"`
	IconURL   string    `json:"iconURL" firestore:"iconURL"`
	IsOnline  bool      `json:"isOnline" firestore:"isOnline"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

// GenerateUUID UUIDを生成する
func GenerateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// InitFirebaseClient Firebaseクライアントを初期化する
func InitFirebaseClient() (*firebase.App, error) {
	// サービスアカウントキーの読み込み
	serviceAccountKey, err := os.ReadFile("config/serviceAccountKey.json")
	if err != nil {
		return nil, fmt.Errorf("サービスアカウントキーの読み込みに失敗: %v", err)
	}

	// サービスアカウントキーのパース
	var serviceAccount struct {
		Type                    string `json:"type"`
		ProjectID               string `json:"project_id"`
		PrivateKeyID            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientID                string `json:"client_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientX509CertURL       string `json:"client_x509_cert_url"`
	}

	if parseErr := json.Unmarshal(serviceAccountKey, &serviceAccount); parseErr != nil {
		return nil, fmt.Errorf("サービスアカウントキーのパースに失敗: %v", parseErr)
	}

	// Firebase初期化オプションの設定
	opt := option.WithCredentialsFile("config/serviceAccountKey.json")

	// Firebaseアプリの初期化
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: serviceAccount.ProjectID,
	}, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase初期化に失敗: %v", err)
	}

	return app, nil
}

// GetUserByEmail メールアドレスでユーザーを検索する
func GetUserByEmail(email string) (*User, error) {
	client, err := InitFirebase()
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

	var user User
	if err := docs[0].DataTo(&user); err != nil {
		return nil, err
	}

	// ドキュメントIDをユーザーIDとして設定
	user.ID = docs[0].Ref.ID

	return &user, nil
}
