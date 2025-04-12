package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
)

// CreateSession セッションを作成する
func CreateSession(user *domain.User) (*domain.Session, error) {
	// セッションIDの生成
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	sessionID := base64.URLEncoding.EncodeToString(bytes)

	// セッションの作成
	session := &domain.Session{
		ID:        sessionID,                           // セッションID
		User:      user,                                // ユーザー
		Token:     sessionID,                           // セッショントークン
		CreatedAt: time.Now(),                          // セッションの作成日時
		UpdatedAt: time.Now(),                          // セッションの更新日時
		ExpiredAt: time.Now().Add(30 * 24 * time.Hour), // 30日間有効
		IsValid:   true,                                // セッションが有効かどうか
	}

	// Firestoreにセッションを保存
	err := firebase.AddData("sessions", session, sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// SetSessionCookie セッションクッキーを設定する
func SetSessionCookie(w http.ResponseWriter, session *domain.Session) {
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
func ValidateSession(w http.ResponseWriter, r *http.Request) (*domain.Session, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("セッションクッキーなし:", err)
		return nil, err
	}

	// Firestoreからセッションを取得
	client, err := firebase.InitFirebase()
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

	var session domain.Session
	if err := doc.DataTo(&session); err != nil {
		fmt.Println("セッションデータ変換エラー:", err)
		return nil, err
	}

	// セッションの有効性をチェック
	if !session.CheckSession() {
		fmt.Println("セッションが無効です")
		return nil, fmt.Errorf("セッションが無効です")
	}

	fmt.Println("")
	return &session, nil
}

// UpdateSession セッションを更新する
func UpdateSession(w http.ResponseWriter, r *http.Request, session *domain.Session) error {
	// Firestoreにセッションを保存（セッションIDをドキュメントIDとして使用）
	err := firebase.AddData("sessions", session, session.ID)
	if err != nil {
		return err
	}

	// セッションクッキーを更新
	SetSessionCookie(w, session)

	return nil
}

// DeleteSession セッションを削除する
func DeleteSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return err
	}

	// Firestoreからセッションを削除
	client, err := firebase.InitFirebase()
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
