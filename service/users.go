// ログインユーザー関連のコントローラー
package service

import (
	"time"
)

type User struct {
	ID            string
	UUID          string
	Name          string
	Email         string
	Password      string // bcryptでハッシュ化されたパスワード
	AvatarURL     string
	StatusMessage string
	IsOnline      bool
	LastLogin     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	IsDeleted     bool
}

type Session struct {
	ID        string
	UUID      string
	Email     string
	UserID    string
	CreatedAt time.Time
}
