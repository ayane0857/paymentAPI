FROM golang:1.20-alpine

RUN mkdir -p /app
# 必要なパッケージをインストール
RUN apk add --no-cache git

WORKDIR /app

# 依存関係ファイルをコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]