# internal ディレクトリについて

この `internal/` ディレクトリは、アプリケーションの内部実装（ビジネスロジックやインフラ層など）を格納するためのディレクトリです。

Goでは、`internal/` 以下のパッケージは **このプロジェクト内からのみインポート可能** であり、他のプロジェクトからのアクセスを防ぐことができます。

クリーンアーキテクチャの思想に基づいて、以下のような層に分かれています。

## ディレクトリ構成
```
internal/
├── config/
│   └── config.go
├── domain/
│   ├── user.go
│   └── product.go
├── usecase/
│   ├── user/
│   │   └── service.go
│   └── product/
│       └── service.go
├── interface/
│   ├── handler/
│   │   ├── home_handler.go
│   │   ├── auth_handler.go
│   │   └── user_handler.go
│   └── firebase/
│       ├── auth.go
│       └── firestore.go
├── infrastructure/
│   ├── firebase/
│   │   └── client.go
│   └── router/
│       └── router.go
```

## 各ディレクトリの役割

### config/
- 初期設定を行うディレクトリ
- グローバル設定（Firebase認証キーなど）もここで管理

### domain/
- ドメインエンティティ（例: `User`, `Product`）を定義します。
- ビジネスルールに直接関わる構造体とインターフェースのみを記述。

### usecase/
- 各ドメインに対応するユースケース（ビジネスロジック）を記述します。
- データの取得・変換・検証など、アプリケーションの中心的な処理を担います。

### interface/
- 外部とのやりとりに関する処理を記述します。
  - `handler/`: HTTP リクエストを処理するハンドラー。
  - `firebase/`: Firebase Auth や Firestore とのやりとりを抽象化。

### infrastructure/
- 具体的な外部技術（Firebase SDKやルーターなど）を実装します。
- `interface/` から呼ばれることが多く、依存関係の末端になります。

---

## 注意事項

- この `internal/` ディレクトリ内のコードは、**他プロジェクトからインポートできません**。
- クリーンアーキテクチャの原則に従い、**内側の層（domain, usecase）ほど外部に依存しない**ように構成してください。
- 新たな機能追加の際は、まず「どの層に位置付けるべきか」を判断してからディレクトリを選びましょう。

---

## 関連資料

- [The Clean Architecture（Uncle Bob）](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

