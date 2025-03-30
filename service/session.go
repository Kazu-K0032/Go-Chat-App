package service

import (
	"crypto/rand"
	"encoding/base64"
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
	}

	// Firestoreにセッションを保存
	err := repository.AddData("sessions", session)
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
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400 * 30, // 30日
	}
	http.SetCookie(w, cookie)
}

// ValidateSession セッションを検証する
func ValidateSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	// Firestoreからセッションを取得
	client, err := repository.InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := r.Context()
	doc, err := client.Collection("sessions").Doc(cookie.Value).Get(ctx)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := doc.DataTo(&session); err != nil {
		return nil, err
	}

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
