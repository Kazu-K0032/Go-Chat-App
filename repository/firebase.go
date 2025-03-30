package repository

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitFirebase() (*firestore.Client, error) {
	opt := option.WithCredentialsFile("config/serviceAccountKey.json")

	// Firebase設定を明示的に指定
	config := &firebase.Config{
		ProjectID: "go-chat-app-cf888",
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Printf("Firebaseアプリの初期化に失敗: %v", err)
		return nil, err
	}

	// タイムアウトを設定したコンテキストを使用
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("Firestoreクライアント作成に失敗: %v", err)
		return nil, err
	}

	return client, nil
}

// コレクションにデータを追加する
func AddData(collection string, data interface{}) error {
	client, err := InitFirebase()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	_, _, err = client.Collection(collection).Add(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
