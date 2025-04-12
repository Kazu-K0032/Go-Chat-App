package router

import (
	"net/http"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/interface/handler"
)

// ルーティングの設定
func SetupRouter(chatUsecase domain.ChatUsecase) *http.ServeMux {
	mux := http.NewServeMux()
	// 静的ファイル (CSS/JS)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("internal/web/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("internal/web/js"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("internal/web/images"))))
	// 静的ファイルの提供
	fs := http.FileServer(http.Dir("internal/web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	// ルーティングの設定
	mux.HandleFunc("/", handler.SearchHandler)
	mux.HandleFunc("/login", handler.LoginHandler)
	mux.HandleFunc("/logout", handler.LogoutHandler)
	mux.HandleFunc("/signup", handler.SignupHandler)
	mux.HandleFunc("/signup/confirm", handler.SignupConfirmHandler)
	mux.HandleFunc("/reset-password", handler.ResetPasswordHandler)
	mux.HandleFunc("/profile", handler.ProfileHandler)
	mux.HandleFunc("/profile/icon", handler.ProfileIconHandler)
	mux.HandleFunc("/chat/", handler.StartChatHandler)
	mux.HandleFunc("/chat", handler.ChatHandler)
	mux.HandleFunc("/search", handler.SearchHandler)
	mux.HandleFunc("/settings", handler.SettingsHandler)

	return mux
}
