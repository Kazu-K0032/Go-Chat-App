package domain

type Session struct {
	Token string
	// その他セッションに必要な情報（例：UserIDなど）
}

func (sess *Session) CheckSession() (bool, error) {
	// 実装例：トークンの検証や有効期限チェックなど
	return true, nil
}
