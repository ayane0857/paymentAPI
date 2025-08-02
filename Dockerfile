FROM golang:1.20-alpine

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