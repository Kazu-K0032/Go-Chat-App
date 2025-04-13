# security_chat_app

Goを使用した基本的なチャットアプリになります。

## スクリーンショット(2925.04.13(日)時点)

## 実装済み機能
- 認証機能（登録/ログイン/ログアウト）
- プロフィール（ユーザー名/画像/パスワードなどの変更）
- 検索機能（登録済みユーザーのフィルタリング）
- チャット機能（他ユーザーと連絡）

## 使用技術
- Go 1.24.1
- Firebase(firestore, storage)
- HTML/CSS(SCSS)/JavaScript
- stylelint, prettier, gulp

## ローカルへの導入手順

1. プロジェクトのclone
  ```:bash
  git clone git@github.com:Kazu-K0032/security_chat_app.git
  ```

2. firebaseからFirebase Admin SDKの認証ファイルを取り込む
  * [Firebase](https://console.firebase.google.com/u/1/?hl=ja)からプロジェクトを作成
  * 作成したプロジェクトにアクセスし、「プロジェクトの設定」⇒「サービスアカウント」⇒「新しい鍵を生成」
    ![security_chat_app_readme_1](https://github.com/user-attachments/assets/c0820422-87d5-4490-80aa-cfe02c564456)
    ![security_chat_app_readme_2](https://github.com/user-attachments/assets/de34f37d-d44b-40a4-8e6f-44ec215f11c9)
  * ダウンロードしたファイル名を「serviceAccountKey.json」に変更し、クローンしたプロジェクトの`internal/config/`に配置してください
  * Firebaseプロジェクト⇒「Firestore Database」から、データベースを作成
  * Firestoreの「ルール」から以下のルールに変更
    
    ```json
    rules_version = '2';
    
    service cloud.firestore {
      match /databases/{database}/documents {
        match /{document=**} {
          allow read, write: if false;
        }
      }
    }
    ```

3. モジュール初期化および依存解決

* 事前に、Go及びNode.jsをダウンロードしてください。
* バージョンは、Goは最低1.21以上, Node.jsはv16.0.0以上
    ```:bash
    cd security_chat_app/
    # Go モジュールの初期化
    go mod tidy
    # Node.jsの依存解決
    npm install
    ```

5. 設定ファイルの確認・修正（`config.ini`）

  * ポートの設定をしています。ご自身の環境に合わせて、随時修正してください。

6. サーバーの起動
   ```:bash
   go run cmd/app/main.go
   ```
   * 実行後、debug.logが生成されます
   * デフォルトだと、`localhost:8050`にアクセスできるようになります。
