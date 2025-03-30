package main

import (
	"log"
	"net/http"

	_ "security_chat_app/config"
	"security_chat_app/repository"
	"security_chat_app/router"
	"security_chat_app/service"
)

func main() {
	client, err := repository.InitFirebase()
	if err != nil {
		log.Fatalf("Firebase初期化に失敗: %v", err)
	}
	defer client.Close()

	chatRepo := repository.NewChatRepository(client)
	chatUsecase := service.NewChatUsecase(chatRepo)

	// ルーティングの設定
	mux := router.SetupRouter(chatUsecase)

	// セッション管理のミドルウェアを適用
	handler := router.Middleware(mux)

	// サーバーの起動
	log.Println("サーバーを起動します...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
