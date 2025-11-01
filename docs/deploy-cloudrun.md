# Google Cloud Run へのデプロイ手順

[English](./lang/en-deploy-cloudrun.md) | 日本語

## 現在の設定（厳密な課金防止設定）

| 項目 | 値 |
|------|-----|
| 最大インスタンス数 | 1 |
| 最小インスタンス数 | 0 |
| CPU | 0.5 vCPU |
| メモリ | 256Mi |
| 同時リクエスト数 | 1 |
| タイムアウト | 60秒 |

**特徴**: 月数人のアクセス程度であれば、無料枠内で運用可能。DoS攻撃による課金を最小限に抑制。

## デプロイコマンド

```bash
PROJECT_ID="xxx"
STORAGE_BUCKET="xxx"
REGION=xxx

gcloud run deploy go-chat-app \
  --image gcr.io/$(gcloud config get-value project)/go-chat-app:latest \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --memory 256Mi \
  --cpu 0.5 \
  --timeout 60 \
  --max-instances 1 \
  --min-instances 0 \
  --concurrency 1 \
  --set-env-vars "PROJECT_ID=${PROJECT_ID},STORAGE_BUCKET=${STORAGE_BUCKET},DEFAULT_ICON_DIR=internal/web/images/defaultIcon,STATIC_DIR=app/views"
```

## 設定確認

```bash
# 現在の設定を確認
gcloud run services describe go-chat-app --region asia-northeast1 \
  --format="table(spec.template.spec.containerConcurrency,spec.template.spec.containers[0].resources.limits,spec.template.metadata.annotations.'autoscaling.knative.dev/maxScale')"

# サービスURLを確認
gcloud run services describe go-chat-app --region asia-northeast1 \
  --format="value(status.url)"
```

## 予算アラート設定（必須）

1. [Google Cloud Console](https://console.cloud.google.com/billing/budgets) にアクセス
2. 「予算を作成」→ 予算額: 1000円、アラート: 50%, 90%, 100%
3. 通知先メールアドレスを設定

予算アラートがないと、異常トラフィック時に課金が発生する可能性があります。

## アップデート手順

```bash
# イメージを再ビルド
docker build -t gcr.io/$(gcloud config get-value project)/go-chat-app:latest .

# イメージをプッシュ
docker push gcr.io/$(gcloud config get-value project)/go-chat-app:latest

# 再デプロイ
gcloud run deploy go-chat-app \
  --image gcr.io/$(gcloud config get-value project)/go-chat-app:latest \
  --region asia-northeast1
```

## トラブルシューティング

```bash
# ログを確認
gcloud run services logs read go-chat-app --region asia-northeast1 --limit 20

# サービスを一時的に無効化（緊急時）
gcloud run services update go-chat-app --region asia-northeast1 --no-traffic
```

## 参考リンク

- [Cloud Run 料金](https://cloud.google.com/run/pricing)
- [予算とアラート](https://console.cloud.google.com/billing/budgets)
