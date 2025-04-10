package domain

import "time"

// type Session struct {
// 	Token string
// 	// その他セッションに必要な情報（例：UserIDなど）
// }

// セッション情報を管理する構造体
type Session struct {
	ID        string
	UserID    string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
	IsValid   bool
}

// CheckSession セッションが有効かチェックする
func (sess *Session) CheckSession() bool {
	// セッションが存在しない場合
	if sess == nil {
		return false
	}

	// 有効期限のチェック
	if time.Now().After(sess.ExpiredAt) {
		return false
	}

	// トークンの検証
	if sess.Token == "" {
		return false
	}

	// IsValidフラグのチェック
	return sess.IsValid
}
