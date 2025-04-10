package router

import (
	"net/http"

	"security_chat_app/internal/config"
	"security_chat_app/internal/usecase/chat"
)

// contextKey コンテキストのキーとして使用するカスタム型
type contextKey string

const templateDataKey contextKey = "templateData"

func StartMainServer(chatUsecase chat.ChatUsecase) error {
	mux := SetupRouter(chatUsecase)
	return http.ListenAndServe(":"+config.Config.Port, mux)
}
