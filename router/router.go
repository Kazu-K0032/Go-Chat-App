package router

// func SetupRouter(chatUsecase service.ChatUsecase) *http.ServeMux {
// 	mux := http.NewServeMux()

// 	// 静的ファイルの提供
// 	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("app/css"))))
// 	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("app/js"))))
// 	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("app/images"))))

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

// 	// その他のリクエストはホームページにリダイレクト
// 	mux.HandleFunc("/", service.HomeHandler)

// 	return mux
// }
