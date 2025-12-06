# FROM golang:1.25-alpine AS builde

# # air をインストール（リポジトリ変更済み）
# RUN go install github.com/air-verse/air@latest

# # ワーキングディレクトリを設定
# WORKDIR /app

# # Goモジュールを取得するためにgo.modとgo.sumをコピー
# COPY go.mod go.sum ./
# ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct 
# RUN go mod download

# # アプリケーションのソースコードをコピー
# COPY . .

# # アプリケーションをビルド
# RUN go build -o main .