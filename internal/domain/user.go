package domain

import "time"

// ユーザーの構造体
type User struct {
	ID            string    // ユーザーのID
	UUID          string    // ユーザーのUUID
	Name          string    // ユーザーの名前
	Email         string    // ユーザーのメールアドレス
	Password      string    // ユーザーのパスワード
	CreatedAt     time.Time // ユーザーの作成日時
	UpdatedAt     time.Time // ユーザーの更新日時
	RememberToken string    // ユーザーのリメンバートークン
	Slug          string    // ユーザーのスラグ
	IsOnline      bool      // ユーザーがオンラインかどうか
	Icon          string    // ユーザーのアイコン
	Gender        string    // ユーザーの性別
}

// ユーザーの関係性の構造体
type Relationship struct {
	ID         string    // 関係性のID
	FollowerID string    // フォロワーのID
	FollowedID string    // フォローされたユーザーのID
	CreatedAt  time.Time // 関係性の作成日時
	UpdatedAt  time.Time // 関係性の更新日時
}
