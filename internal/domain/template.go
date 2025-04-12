package domain

// TemplateData 共通のテンプレートデータ構造体
type TemplateData struct {
	IsLoggedIn       bool       // ログイン状態
	User             *User      // ユーザー情報
	Messages         []Message  // メッセージ
	Contacts         []Contact  // 連絡先
	Chats            []Chat     // チャット
	CurrentChat      *Chat      // 現在のチャット
	SignupForm       SignupForm // サインアップフォーム
	LoginForm        LoginForm  // ログインフォーム
	ResetForm        ResetForm  // リセットフォーム
	ValidationErrors []string   // バリデーションエラー
	Error            string     // エラー
}

// DefaultIcon デフォルトアイコンの情報
type DefaultIcon struct {
	Path string // パス
	Name string // 名前
}

// GetDefaultIcons デフォルトアイコンの一覧を取得する
func GetDefaultIcons() []DefaultIcon {
	return []DefaultIcon{}
}
