# ビルドステージ
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 依存パッケージのインストール
RUN apk add --no-cache git

# Goモジュールのコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# バイナリのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .
COPY --from=builder /app/config.ini .
COPY --from=builder /app/internal/web ./internal/web

# Cloud Runは環境変数PORTを使用するので、デフォルトは8080
ENV PORT=8080

EXPOSE 8080

CMD ["./main"]

