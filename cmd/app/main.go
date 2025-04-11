package main

import (
	"log"

	"security_chat_app/internal/infrastructure/firebase"
)

func main() {
	client, err := firebase.InitFirebase()
	if err != nil {
		log.Fatalf("Firebase初期化に失敗: %v", err)
	}
	defer client.Close()

	// チャットリポジトリの作成
	chatRepo := repository.NewChatRepository(client)
	chatUsecase := chat.NewChatUsecase(chatRepo)
	if chatUsecase == nil {
		log.Fatal("chatUsecaseがnilです")
	}

	// ルーティングの設定
	// mux := router.SetupRouter(nil) // chatUsecaseが未実装の場合はnilを渡す

	// サーバーを起動
	// log.Printf("サーバーを起動します。ポート: %s", config.Config.Port)
	// if err := http.ListenAndServe(":"+config.Config.Port, mux); err != nil {
	// 	log.Fatal(err)
	// }
}
