package router

// func SetupRouter(chatUsecase service.ChatUsecase) *http.ServeMux {
// 	mux := http.NewServeMux()

// 	// 静的ファイルの提供
// 	fs := http.FileServer(http.Dir("app"))
// 	mux.Handle("/", fs)

// 	// 認証関連
// 	mux.HandleFunc("/signup", service.SignupHandler)
// 	mux.HandleFunc("/login", service.LoginHandler)
// 	mux.HandleFunc("/logout", service.LogoutHandler)

// 	// プロフィール関連
// 	mux.HandleFunc("/profile", service.ProfileHandler)
// 	mux.HandleFunc("/profile/icon", service.ProfileIconHandler)

// 	// 設定関連
// 	mux.HandleFunc("/settings", service.SettingsHandler)

// 	// チャット関連
// 	mux.HandleFunc("/chat", service.ChatHandler)
// 	mux.HandleFunc("/chat/start/", service.StartChatHandler)
// 	mux.HandleFunc("/chat/messages/", func(w http.ResponseWriter, r *http.Request) {
// 		service.ChatHandler(w, r)
// 	})

// 	return mux
// }
