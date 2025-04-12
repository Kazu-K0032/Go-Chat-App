package domain

import "time"

// セッション情報を管理する構造体
type Session struct {
	ID        string    // セッションのID
	User      *User     // ユーザー
	Token     string    // セッションのトークン
	CreatedAt time.Time // セッションの作成日時
	UpdatedAt time.Time // セッションの更新日時
	ExpiredAt time.Time // セッションの有効期限
	IsValid   bool      // セッションが有効かどうか
}

// Sessionの有効性をチェックする
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
