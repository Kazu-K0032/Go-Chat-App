package domain

// サインアップフォームのデータ構造体
type SignupForm struct {
	Name     string
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
	Email           string
	Password        string
	PasswordConfirm string
}

// パスワード変更フォームのデータ構造体
type PasswordForm struct {
	CurrentPassword    string
	NewPassword       string
	NewPasswordConfirm string
}
