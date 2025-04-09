package repository

import "context"

// GetUserByEmail メールアドレスでユーザーを検索する
func GetUserByEmail(email string) (*User, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	query := client.Collection("users").Where("email", "==", email)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil // ユーザーが見つからない
	}

	var user User
	if err := docs[0].DataTo(&user); err != nil {
		return nil, err
	}

	// ドキュメントIDをユーザーIDとして設定
	user.ID = docs[0].Ref.ID

	return &user, nil
}
