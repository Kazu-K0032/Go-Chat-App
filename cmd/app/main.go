package main

import (
	"log"
	"net/http"

	"security_chat_app/internal/config"
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

	// サーバーを起動
	log.Printf("サーバーを起動します。ポート: %s", config.Config.Port)
	if err := http.ListenAndServe(":"+config.Config.Port, handler); err != nil {
		log.Fatal(err)
	}
}
