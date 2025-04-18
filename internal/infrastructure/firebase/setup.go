package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"security_chat_app/internal/config"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func InitFirebase() (*firestore.Client, error) {
	opt := option.WithCredentialsFile(config.Config.ServiceKeyPath)

	// Firebase設定を明示的に指定
	firebaseConfig := &firebase.Config{
		ProjectID:     config.Config.ProjectId,
		StorageBucket: config.Config.StorageBucket,
	}

	app, err := firebase.NewApp(context.Background(), firebaseConfig, opt)
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

	// デフォルトアイコンの初期化
	if err := initDefaultIcons(app); err != nil {
		log.Printf("デフォルトアイコンの初期化に失敗: %v", err)
	}

	return client, nil
}

// デフォルトアイコンを初期化する
func initDefaultIcons(app *firebase.App) error {
	ctx := context.Background()
	
	// Storageクライアントを取得
	storageClient, err := app.Storage(ctx)
	if err != nil {
		return fmt.Errorf("Storageクライアントの作成に失敗: %v", err)
	}
	
	// バケットを取得
	bucket, err := storageClient.DefaultBucket()
	if err != nil {
		return fmt.Errorf("デフォルトバケットの取得に失敗: %v", err)
	}
	
	// デフォルトアイコンディレクトリの存在確認
	prefix := config.Config.DefaultIconDir
	it := bucket.Objects(ctx, &storage.Query{Prefix: prefix})
	
	// 少なくとも1つのオブジェクトが存在するか確認
	hasObjects := false
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("オブジェクトの列挙に失敗: %v", err)
		}
		hasObjects = true
		break
	}
	
	// デフォルトアイコンが存在しない場合、作成する
	if !hasObjects {
		// ディレクトリ内のファイルを取得
		files, err := os.ReadDir(config.Config.DefaultIconDir)
		if err != nil {
			return fmt.Errorf("デフォルトアイコンディレクトリの読み込みに失敗: %v", err)
		}
		// 各ファイルをアップロード
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			
			// ファイルパス
			filePath := filepath.Join(config.Config.DefaultIconDir, file.Name())
			
			// ファイルを開く
			fileContent, err := os.Open(filePath)
			if err != nil {
				log.Printf("ファイル %s のオープンに失敗: %v", filePath, err)
				continue
			}
			defer fileContent.Close()
			
			// アップロード先のパス
			objectPath := prefix + file.Name()
			
			// オブジェクトを作成
			obj := bucket.Object(objectPath)
			writer := obj.NewWriter(ctx)
			
			// メタデータを設定
			writer.ObjectAttrs = storage.ObjectAttrs{
				Name:        objectPath,
				ContentType: getContentType(file.Name()),
				ACL:         []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}},
			}
			
			// ファイルをアップロード
			if _, err := io.Copy(writer, fileContent); err != nil {
				log.Printf("ファイル %s のアップロードに失敗: %v", filePath, err)
				writer.Close()
				continue
			}
			
			// ライターを閉じる
			if err := writer.Close(); err != nil {
				log.Printf("ファイル %s のアップロード完了に失敗: %v", filePath, err)
				continue
			}
			
			log.Printf("デフォルトアイコン %s をアップロードしました", file.Name())
		}
	}
	
	return nil
}

// ファイル名からContentTypeを取得
func getContentType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// InitFirebaseClient Firebaseクライアントを初期化する
func InitFirebaseClient() (*firebase.App, error) {
	// サービスアカウントキーの読み込み
	serviceAccountKey, err := os.ReadFile(config.Config.ServiceKeyPath)
	if err != nil {
		return nil, fmt.Errorf("サービスアカウントキーの読み込みに失敗: %v", err)
	}

	// サービスアカウントキーのパース
	var serviceAccount struct {
		Type                    string `json:"type"`
		ProjectID               string `json:"project_id"`
		PrivateKeyID            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientID                string `json:"client_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientX509CertURL       string `json:"client_x509_cert_url"`
	}

	if parseErr := json.Unmarshal(serviceAccountKey, &serviceAccount); parseErr != nil {
		return nil, fmt.Errorf("サービスアカウントキーのパースに失敗: %v", parseErr)
	}

	// Firebase初期化オプションの設定
	opt := option.WithCredentialsFile(config.Config.ServiceKeyPath)

	// Firebaseアプリの初期化
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID:     serviceAccount.ProjectID,
		StorageBucket: config.Config.StorageBucket, // Storageバケット名を明示的に指定
	}, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase初期化に失敗: %v", err)
	}

	return app, nil
}
