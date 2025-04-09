# ディレクトリ構成

```
project-root/
├── cmd/
│   └── app/
│       └── main.go          // アプリのエントリーポイント
├── internal/
│   ├── domain/              // エンティティやドメインモデル
│   │   └── user.go
│
│   ├── usecase/             // ユースケース層（ビジネスロジック）
│   │   └── user/
│   │       └── service.go
│
│   ├── interface/           // インターフェース層（外部との接続）
│   │   ├── handler/         // HTTPハンドラー
│   │   │   ├── home_handler.go
│   │   │   ├── auth_handler.go
│   │   │   └── user_handler.go
│   │   └── firebase/        // Firebase SDK操作のラッパー
│   │       ├── auth.go
│   │       └── firestore.go
│
│   ├── infrastructure/      // データベースや外部APIとのやり取り
│   │   ├── firebase/
│   │   │   └── client.go    // Firebaseの初期化や設定
│   │   └── router/
│   │       └── router.go    // ルーティング定義（Ginなど）
│
│   └── config/
│       └── config.go        // 設定の読み込み（envファイルなど）
│
├── web/                     // テンプレートや静的ファイル（必要であれば）
│   ├── templates/
│   └── static/
│
├── test/                    // 単体・統合テスト
│   └── handler_test.go
│
├── go.mod
└── go.sum
```
