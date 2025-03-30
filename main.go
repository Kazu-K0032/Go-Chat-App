package main

import (
	"log"

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

	if err := router.StartMainServer(chatUsecase); err != nil {
		log.Fatalf("サーバー起動に失敗: %v", err)
	}
}
