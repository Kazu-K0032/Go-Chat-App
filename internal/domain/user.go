package domain

import "time"

// type User struct {
// 	ID            string
// 	UUID          string
// 	Name          string
// 	Email         string
// 	Password      string // bcryptでハッシュ化されたパスワード
// 	AvatarURL     string
// 	StatusMessage string
// 	IsOnline      bool
// 	LastLogin     time.Time
// 	CreatedAt     time.Time
// 	UpdatedAt     time.Time
// 	IsDeleted     bool
// }

type User struct {
	ID             string
	UUID           string
	Name           string
	Email          string
	PasswordDigest string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RememberToken  string
	Slug           string
	IsOnline       bool
	Icon           string
}

type Relationship struct {
	ID         string
	FollowerID string
	FollowedID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
