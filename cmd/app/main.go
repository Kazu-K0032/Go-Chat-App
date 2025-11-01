package main

import (
	"log"
	"net/http"
	"os"

	"security_chat_app/internal/config"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/infrastructure/router"
	"security_chat_app/internal/usecase/chat"
)

func main() {
	// Firebaseの初期化
	client, err := firebase.InitFirebase()
	if err != nil {
		log.Fatalf("Firebase初期化に失敗: %v", err)
	}
	defer client.Close()

	// チャットリポジトリの作成
	chatRepo := chat.NewChatRepository(client)
	chatUsecase := chat.NewChatUsecase(chatRepo)
	if chatUsecase == nil {
		log.Fatal("チャットのユースケースの実装に不備があります")
	}

	// ルーティングの設定
	httpRouter := router.SetupRouter(chatUsecase)
	if httpRouter == nil {
		log.Fatal("ルーティングの設定に不備があります")
	}

	// サーバーを起動
	// Cloud Runは環境変数PORTを自動設定するため、それを優先する
	port := config.Config.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
		log.Printf("環境変数PORTを検出しました: %s", envPort)
	}
	log.Printf("サーバーを起動します。ポート: %s", port)
	if err := http.ListenAndServe(":"+port, httpRouter); err != nil {
		log.Fatal("サーバーの起動に失敗しました")
	}
}
