package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"security_chat_app/internal/config"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Post 投稿情報を持つ構造体
type Post struct {
	ID        string    `json:"id" firestore:"id"`
	Content   string    `json:"content" firestore:"content"`
	UserID    string    `json:"user_id" firestore:"user_id"`
	ReplyTo   string    `json:"reply_to,omitempty" firestore:"reply_to,omitempty"`
	LikedBy   []string  `json:"liked_by" firestore:"liked_by"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}

func InitFirebase() (*firestore.Client, error) {
	opt := option.WithCredentialsFile(config.Config.FirebaseServiceAccountKey)

	// Firebase設定を明示的に指定
	config := &firebase.Config{
		ProjectID: "go-chat-app-cf888",
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Printf("Firebaseアプリの初期化に失敗: %v", err)
		return nil, err
	}

	// タイムアウトを設定したコンテキストを使用
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("Firestoreクライアント作成に失敗: %v", err)
		return nil, err
	}

	return client, nil
}

// コレクションにデータを追加する
func AddData(collection string, data interface{}, docID string) error {
	client, err := InitFirebase()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	if docID != "" {
		// カスタムドキュメントIDを使用
		_, err = client.Collection(collection).Doc(docID).Set(ctx, data)
	} else {
		// 自動生成のドキュメントIDを使用
		_, _, err = client.Collection(collection).Add(ctx, data)
	}
	if err != nil {
		return err
	}
	return nil
}

// コレクションとドキュメントIDから特定フィールドを更新する
func UpdateField(collection string, documentID string, field string, value interface{}) error {
	client, err := InitFirebase()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Collection(collection).Doc(documentID).Update(ctx, []firestore.Update{
		{
			Path:  field,
			Value: value,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// コレクションからデータを取得する
func GetData(collection string, documentID string) (map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	doc, err := client.Collection(collection).Doc(documentID).Get(ctx)
	if err != nil {
		return nil, err
	}

	return doc.Data(), nil
}

// コレクションから条件に合うデータを取得する
func GetDataByQuery(collection string, field string, operator string, value interface{}) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	query := client.Collection(collection).Where(field, operator, value)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, doc := range docs {
		results = append(results, doc.Data())
	}

	return results, nil
}

// コレクションからデータを削除する
func DeleteData(collection string, documentID string) error {
	client, err := InitFirebase()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Collection(collection).Doc(documentID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

// コレクションの全データを取得する
func GetAllData(collection string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	docs, err := client.Collection(collection).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID // ドキュメントIDをidフィールドとして追加
		results = append(results, data)
	}

	return results, nil
}

// ユーザーを検索する
func SearchUsers(searchQuery string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()

	// ユーザー名で部分一致検索
	usersQuery := client.Collection("users").
		Where("name", ">=", searchQuery).
		Where("name", "<=", searchQuery+"\uf8ff").
		Limit(20)

	docs, err := usersQuery.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID // ドキュメントIDをidフィールドとして追加
		results = append(results, data)
	}

	return results, nil
}

// チャットを開始する
func StartChat(userID string, targetUserID string) (string, error) {
	client, err := InitFirebase()
	if err != nil {
		return "", err
	}
	defer client.Close()

	ctx := context.Background()

	// チャットIDを生成
	chatID := fmt.Sprintf("chat_%d", time.Now().UnixNano())

	// チャットを作成
	chat := map[string]interface{}{
		"id":           chatID,
		"participants": []string{userID, targetUserID},
		"createdAt":    time.Now(),
		"updatedAt":    time.Now(),
	}

	_, err = client.Collection("chats").Doc(chatID).Set(ctx, chat)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

// チャットメッセージを追加する
func AddChatMessage(chatID string, message map[string]interface{}) error {
	client, err := InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
		return err
	}
	defer client.Close()

	ctx := context.Background()

	// メッセージIDを生成
	messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())
	message["id"] = messageID

	// メッセージを保存
	_, err = client.Collection("chats").Doc(chatID).Collection("messages").Doc(messageID).Set(ctx, message)
	if err != nil {
		log.Printf("メッセージ保存エラー: %v", err)
		return err
	}

	// チャットの更新時刻を更新
	_, err = client.Collection("chats").Doc(chatID).Update(ctx, []firestore.Update{
		{
			Path:  "updated_at",
			Value: time.Now(),
		},
	})
	if err != nil {
		log.Printf("チャット更新時刻の更新エラー: %v", err)
		return err
	}

	log.Printf("メッセージを保存しました: chatID=%s, messageID=%s", chatID, messageID)
	return nil
}

// チャットのメッセージを取得する
func GetChatMessages(chatID string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	docs, err := client.Collection("chats").Doc(chatID).Collection("messages").OrderBy("created_at", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var messages []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID
		messages = append(messages, data)
	}

	return messages, nil
}

// チャットの存在確認
func CheckChatExists(chatID string) (bool, error) {
	client, err := InitFirebase()
	if err != nil {
		return false, err
	}
	defer client.Close()

	ctx := context.Background()
	doc, err := client.Collection("chats").Doc(chatID).Get(ctx)
	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}

// チャットの参加者を取得
func GetChatParticipants(chatID string) ([]string, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	doc, err := client.Collection("chats").Doc(chatID).Get(ctx)
	if err != nil {
		return nil, err
	}

	data := doc.Data()
	participants, ok := data["participants"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("participants field is invalid")
	}

	var result []string
	for _, p := range participants {
		if str, ok := p.(string); ok {
			result = append(result, str)
		}
	}

	return result, nil
}

// アイコンをアップロードする
func UploadIcon(userID string, filePath string) (string, error) {
	// Firebase Storageクライアントを初期化
	opt := option.WithCredentialsFile("config/serviceAccountKey.json")
	config := &firebase.Config{
		ProjectID:     "go-chat-app-cf888",
		StorageBucket: "go-chat-app-cf888.firebasestorage.app",
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return "", fmt.Errorf("firebaseアプリの初期化に失敗しました: %v", err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return "", fmt.Errorf("storageクライアントの作成に失敗しました: %v", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("バケットの取得に失敗しました: %v", err)
	}

	// アップロード先のパスを設定
	objectPath := fmt.Sprintf("icons/%s%s", userID, filepath.Ext(filePath))
	object := bucket.Object(objectPath)

	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("ファイルのオープンに失敗しました: %v", err)
	}
	defer file.Close()

	// ファイルをアップロード
	wc := object.NewWriter(context.Background())

	// メタデータを設定
	wc.ObjectAttrs = storage.ObjectAttrs{
		Name:        objectPath,
		ContentType: "image/jpeg",
		ACL:         []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}},
	}

	_, err = io.Copy(wc, file)
	if err != nil {
		return "", fmt.Errorf("ファイルのアップロードに失敗しました: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("ライターのクローズに失敗しました: %v", err)
	}

	// 公開URLを取得
	attrs, err := object.Attrs(context.Background())
	if err != nil {
		return "", fmt.Errorf("オブジェクトの属性取得に失敗しました: %v", err)
	}

	return attrs.MediaLink, nil
}

// デフォルトアイコンのURLを取得
func GetDefaultIconURL(objectPath string) (string, error) {
	fmt.Printf("デフォルトアイコンを取得中: %s\n", objectPath)

	// 公開URLを生成
	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/go-chat-app-cf888.firebasestorage.app/o/%s?alt=media", url.PathEscape(objectPath))
	fmt.Printf("デフォルトアイコンのURLを取得しました: %s\n", url)

	return url, nil
}

// GetUserPosts ユーザーの投稿を取得する
func GetUserPosts(userID string) ([]Post, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	query := client.Collection("posts").
		Where("user_id", "==", userID).
		Where("reply_to", "==", "").
		OrderBy("created_at", firestore.Desc)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, doc := range docs {
		var post Post
		if err := doc.DataTo(&post); err != nil {
			continue
		}
		post.ID = doc.Ref.ID
		posts = append(posts, post)
	}

	return posts, nil
}

// GetUserReplies ユーザーの返信を取得する
func GetUserReplies(userID string) ([]Post, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	query := client.Collection("posts").
		Where("user_id", "==", userID).
		Where("reply_to", "!=", "").
		OrderBy("created_at", firestore.Desc)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var replies []Post
	for _, doc := range docs {
		var reply Post
		if err := doc.DataTo(&reply); err != nil {
			continue
		}
		reply.ID = doc.Ref.ID
		replies = append(replies, reply)
	}

	return replies, nil
}

// GetUserLikedPosts ユーザーがいいねした投稿を取得する
func GetUserLikedPosts(userID string) ([]Post, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	query := client.Collection("posts").
		Where("liked_by", "array-contains", userID).
		OrderBy("created_at", firestore.Desc)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var likes []Post
	for _, doc := range docs {
		var post Post
		if err := doc.DataTo(&post); err != nil {
			continue
		}
		post.ID = doc.Ref.ID
		likes = append(likes, post)
	}

	return likes, nil
}
