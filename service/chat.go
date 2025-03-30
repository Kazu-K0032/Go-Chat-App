package service

import (
	"net/http"
)

type ChatController struct {
	chatUsecase ChatUsecase
}

func NewChatController(chatUsecase ChatUsecase) *ChatController {
	return &ChatController{
		chatUsecase: chatUsecase,
	}
}

// チャットの作成
func (c *ChatController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.FormValue("user")
	message := r.FormValue("message")

	if err := c.chatUsecase.CreateChat(user, message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
