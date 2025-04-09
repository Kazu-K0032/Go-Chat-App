package handler

import (
	"net/http"

	"security_chat_app/service"
)

// LogoutHandler ログアウト処理を実行
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := service.DeleteSession(w, r)
		if err != nil {
			http.Error(w, "ログアウトエラー", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
