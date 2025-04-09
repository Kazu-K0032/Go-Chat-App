package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"security_chat_app/repository"
)

// Session セッション情報
type Session struct {
	ID        string
	UserID    string
	Email     string
	CreatedAt time.Time
	User      *repository.User
}

// CreateSession セッションを作成する
func CreateSession(user *repository.User) (*Session, error) {
	// セッションIDの生成
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	sessionID := base64.URLEncoding.EncodeToString(bytes)

	session := &Session{
		ID:        sessionID,
		UserID:    user.ID,
		Email:     user.Email,
		CreatedAt: time.Now(),
		User:      user,
	}

	// Firestoreにセッションを保存（セッションIDをドキュメントIDとして使用）
	err := repository.AddData("sessions", session, sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// SetSessionCookie セッションクッキーを設定する
func SetSessionCookie(w http.ResponseWriter, session *Session) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                // 開発環境ではfalseに設定
		SameSite: http.SameSiteLaxMode, // 開発環境ではLaxに設定
		MaxAge:   86400 * 30,           // 30日
	}
	http.SetCookie(w, cookie)
}

// ValidateSession セッションを検証する
func ValidateSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	fmt.Println("セッション検証開始")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("セッションクッキーなし:", err)
		return nil, err
	}

	fmt.Println("セッションクッキー取得:", cookie.Value)
	// Firestoreからセッションを取得
	client, err := repository.InitFirebase()
	if err != nil {
		fmt.Println("Firebase初期化エラー:", err)
		return nil, err
	}
	defer client.Close()

	ctx := r.Context()
	doc, err := client.Collection("sessions").Doc(cookie.Value).Get(ctx)
	if err != nil {
		fmt.Println("セッション取得エラー:", err)
		return nil, err
	}

	var session Session
	if err := doc.DataTo(&session); err != nil {
		fmt.Println("セッションデータ変換エラー:", err)
		return nil, err
	}

	// ユーザー情報を取得
	userData, err := repository.GetData("users", session.UserID)
	if err != nil {
		fmt.Println("ユーザー情報取得エラー:", err)
		return nil, err
	}

	// ユーザー情報をセッションに設定
	session.User = &repository.User{
		ID:        userData["id"].(string),
		Name:      userData["name"].(string),
		Email:     userData["email"].(string),
		Password:  userData["password"].(string),
		CreatedAt: userData["created_at"].(time.Time),
		UpdatedAt: userData["updated_at"].(time.Time),
	}

	fmt.Println("セッション検証成功")
	return &session, nil
}

// DeleteSession セッションを削除する
func DeleteSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return err
	}

	// Firestoreからセッションを削除
	client, err := repository.InitFirebase()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := r.Context()
	_, err = client.Collection("sessions").Doc(cookie.Value).Delete(ctx)
	if err != nil {
		return err
	}

	// クッキーを削除
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	return nil
}

// ログアウト処理
func Logout(w http.ResponseWriter, r *http.Request) {
	// セッションを削除
	err := DeleteSession(w, r)
	if err != nil {
		// エラーが発生してもログアウトは成功したものとして扱う
	}

	// フラッシュメッセージを設定
	flash := &http.Cookie{
		Name:   "flash",
		Value:  "ログアウトしました",
		Path:   "/",
		MaxAge: 1,
	}
	http.SetCookie(w, flash)

	// トップページにリダイレクト
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
