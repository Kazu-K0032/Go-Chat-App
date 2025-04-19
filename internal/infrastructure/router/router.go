package router

import (
	"net/http"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/interface/handler"
)

// ルーティングの設定
func SetupRouter(chatUsecase domain.ChatUsecase) *http.ServeMux {
	rootDir := "internal/web/"
	httpRouter := http.NewServeMux()
	// 静的ファイル (CSS/JS)
	httpRouter.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(rootDir+"css"))))
	httpRouter.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(rootDir+"js"))))
	httpRouter.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(rootDir+"images"))))
	// ルーティング
	httpRouter.HandleFunc("/", handler.SearchHandler)
	httpRouter.HandleFunc("/login", handler.LoginHandler)
	httpRouter.HandleFunc("/logout", handler.LogoutHandler)
	httpRouter.HandleFunc("/signup", handler.SignupHandler)
	httpRouter.HandleFunc("/signup/confirm", handler.SignupConfirmHandler)
	httpRouter.HandleFunc("/reset-password", handler.ResetPasswordHandler)
	httpRouter.HandleFunc("/profile", handler.ProfileHandler)
	httpRouter.HandleFunc("/profile/", handler.ProfileHandler)
	httpRouter.HandleFunc("/profile/icon", handler.ProfileIconHandler)
	httpRouter.HandleFunc("/chat/", handler.StartChatHandler)
	httpRouter.HandleFunc("/chat", handler.ChatHandler)
	httpRouter.HandleFunc("/search", handler.SearchHandler)
	httpRouter.HandleFunc("/settings", handler.SettingsHandler)
	httpRouter.HandleFunc("/settings/username", handler.SettingsHandler)

	return httpRouter
}
