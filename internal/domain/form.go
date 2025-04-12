package domain

// サインアップフォームのデータ構造体
type SignupForm struct {
	Username string
	Email    string
	Password string
}

// ログインフォームのデータ構造体
type LoginForm struct {
	Email    string
	Password string
}

// パスワードリセットフォームのデータ構造体
type ResetForm struct {
	Email string
}
