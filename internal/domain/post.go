package domain

import "time"

// Post 投稿情報を持つ構造体
type Post struct {
	ID        string    `json:"id" firestore:"id"`
	Content   string    `json:"content" firestore:"content"`
	UserID    string    `json:"user_id" firestore:"user_id"`
	ReplyTo   string    `json:"reply_to,omitempty" firestore:"reply_to,omitempty"`
	LikedBy   []string  `json:"liked_by" firestore:"liked_by"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}

type sPosts struct {
	ID        string
	Content   string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	ReplyToID string
	LikedBy   []string
}

type sReplies struct {
	ID        string
	Content   string
	PostID    string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type sLikes struct {
	ID        string
	PostID    string
	UserID    string
	CreatedAt time.Time
}
