# Deployment Guide for Google Cloud Run

[日本語](../deploy-cloudrun.md) | English

## Current Configuration (Strict Cost Prevention Settings)

| Item | Value |
|------|-------|
| Max Instances | 1 |
| Min Instances | 0 |
| CPU | 0.5 vCPU |
| Memory | 256Mi |
| Concurrent Requests | 1 |
| Timeout | 60 seconds |

**Features**: Suitable for a few users per month, can operate within the free tier. Minimizes costs from DoS attacks.

## Deployment Command

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

## Configuration Check

```bash
# Check current configuration
gcloud run services describe go-chat-app --region asia-northeast1 \
  --format="table(spec.template.spec.containerConcurrency,spec.template.spec.containers[0].resources.limits,spec.template.metadata.annotations.'autoscaling.knative.dev/maxScale')"

# Check service URL
gcloud run services describe go-chat-app --region asia-northeast1 \
  --format="value(status.url)"
```

## Budget Alert Settings (Required)

1. Access [Google Cloud Console](https://console.cloud.google.com/billing/budgets)
2. Click "Create Budget" → Budget amount: ¥1,000, Alerts: 50%, 90%, 100%
3. Set notification email address

Without budget alerts, unexpected charges may occur if there is abnormal traffic.

## Update Procedure

```bash
# Rebuild image
docker build -t gcr.io/$(gcloud config get-value project)/go-chat-app:latest .

# Push image
docker push gcr.io/$(gcloud config get-value project)/go-chat-app:latest

# Redeploy
gcloud run deploy go-chat-app \
  --image gcr.io/$(gcloud config get-value project)/go-chat-app:latest \
  --region asia-northeast1
```

## Troubleshooting

```bash
# Check logs
gcloud run services logs read go-chat-app --region asia-northeast1 --limit 20

# Temporarily disable service (emergency)
gcloud run services update go-chat-app --region asia-northeast1 --no-traffic
```

## Reference Links

- [Cloud Run Pricing](https://cloud.google.com/run/pricing)
- [Budgets and Alerts](https://console.cloud.google.com/billing/budgets)

