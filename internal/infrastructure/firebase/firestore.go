package firebase

import (
	"context"
	"fmt"
	"log"
	"time"

	"security_chat_app/internal/domain"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// コレクションにデータを追加する
func AddData(collection string, data interface{}, docID string) error {
	client, err := InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
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
		log.Printf("データ追加エラー: %v", err)
		return err
	}
	log.Printf("データを追加しました: collection=%s, docID=%s", collection, docID)
	return nil
}

// コレクションとドキュメントIDから特定フィールドを更新する
func UpdateField(collection string, documentID string, field string, value interface{}) error {
	client, err := InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
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
		log.Printf("フィールド更新エラー: %v, collection=%s, documentID=%s, field=%s", err, collection, documentID, field)
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
func GetAllData(collection string, userID string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	
	var docs []*firestore.DocumentSnapshot
	var err2 error
	
	if collection == "chats" {
		// チャットの場合は、participantsフィールドに基づいてフィルタリング
		query := client.Collection(collection).Where("participants", "array-contains", userID)
		docs, err2 = query.Documents(ctx).GetAll()
	} else {
		// その他のコレクションは従来通り全件取得
		docs, err2 = client.Collection(collection).Documents(ctx).GetAll()
	}
	
	if err2 != nil {
		return nil, err2
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
func SearchUser(searchQuery string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()

	// ユーザー名で部分一致検索（大文字小文字を区別しない）
	usersQuery := client.Collection("users").
		Where("Name", ">=", searchQuery).
		Where("Name", "<=", searchQuery+"\uf8ff").
		Limit(20)

	docs, err := usersQuery.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		return nil, err
	}

	var results []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		data["ID"] = doc.Ref.ID // ドキュメントIDをIDフィールドとして追加
		log.Printf("検索結果のユーザーデータ: %+v", data)
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

// GetUserPosts ユーザーの投稿を取得する
func GetUserPosts(userID string) ([]domain.Post, error) {
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

	var posts []domain.Post
	for _, doc := range docs {
		var post domain.Post
		if err := doc.DataTo(&post); err != nil {
			continue
		}
		post.ID = doc.Ref.ID
		posts = append(posts, post)
	}

	return posts, nil
}

// GetUserReplies ユーザーの返信を取得する
func GetUserReplies(userID string) ([]domain.Post, error) {
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

	var replies []domain.Post
	for _, doc := range docs {
		var reply domain.Post
		if err := doc.DataTo(&reply); err != nil {
			continue
		}
		reply.ID = doc.Ref.ID
		replies = append(replies, reply)
	}

	return replies, nil
}

// GetUserLikedPosts ユーザーがいいねした投稿を取得する
func GetUserLikedPosts(userID string) ([]domain.Post, error) {
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

	var likes []domain.Post
	for _, doc := range docs {
		var post domain.Post
		if err := doc.DataTo(&post); err != nil {
			continue
		}
		post.ID = doc.Ref.ID
		likes = append(likes, post)
	}

	return likes, nil
}

// GetAllChats は指定されたユーザーIDが参加者として含まれるチャットを全て取得します
func GetAllChats(userID string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()

	// チャットコレクションを参照
	chatsRef := client.Collection("chats")

	// ユーザーIDが参加者に含まれるチャットを検索
	query := chatsRef.Where("participants", "array-contains", userID)
	iter := query.Documents(ctx)
	defer iter.Stop()

	var chats []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("チャットデータの取得に失敗: %v", err)
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		chats = append(chats, data)
	}

	return chats, nil
}
