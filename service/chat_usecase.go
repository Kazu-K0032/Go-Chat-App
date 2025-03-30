package service

type ChatUsecase interface {
	CreateChat(user, message string) error
}

type chatUsecase struct {
	repo interface {
		AddChat(user, message string) error
	}
}

func NewChatUsecase(repo interface {
	AddChat(user, message string) error
},
) ChatUsecase {
	return &chatUsecase{repo: repo}
}

// チャット作成時のビジネスロジックを定義
func (c *chatUsecase) CreateChat(user, message string) error {
	return c.repo.AddChat(user, message)
}
