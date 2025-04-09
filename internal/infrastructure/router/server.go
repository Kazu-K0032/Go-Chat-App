package router

import (
	"net/http"

	"security_chat_app/service"
)

// contextKey コンテキストのキーとして使用するカスタム型
type contextKey string

const templateDataKey contextKey = "templateData"

func StartMainServer(chatUsecase service.ChatUsecase) error {
	mux := SetupRouter(chatUsecase)
	return http.ListenAndServe(":8050", mux)
}
