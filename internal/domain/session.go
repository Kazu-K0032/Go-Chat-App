package domain

import (
	"log"
	"time"
)

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

// CheckSession セッションの有効性をチェックする
func (s *Session) CheckSession() bool {
	log.Printf("セッション有効性チェック開始: sessionID=%s", s.ID)

	if s == nil {
		log.Printf("セッションがnilです")
		return false
	}

	// セッションの有効期限をチェック
	if time.Now().After(s.ExpiredAt) {
		log.Printf("セッションの有効期限が切れています: sessionID=%s, expiredAt=%v", s.ID, s.ExpiredAt)
		return false
	}

	// セッションが無効に設定されている場合
	if !s.IsValid {
		log.Printf("セッションが無効に設定されています: sessionID=%s", s.ID)
		return false
	}

	log.Printf("セッション有効性チェック成功: sessionID=%s", s.ID)
	return true
}
