FROM golang:1.24.5-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git

# 作業ディレクトリ作成
WORKDIR /app

# go.mod, go.sum をコピーして依存解決
COPY go.mod go.sum ./
# Go Modules Proxyを明示
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN go build -o main .

# 実行用の軽量イメージ
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]