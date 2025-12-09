package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	pkg "kitta_backend/pkg"
	"kitta_backend/tutorial"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

// server は gRPC サーバー構造体
type server struct {
	pkg.UnimplementedCutListServiceServer
	queries *tutorial.Queries // DB 用クエリを保持
}

// GetUserByUuid は gRPC メソッドの実装
func (s *server) GetUserByUuid(ctx context.Context, in *pkg.ReqGetUserByUuid) (*pkg.ResGetUserByUuid, error) {
	// SQLC の GetUser を使って ID を取得
	id, err := s.queries.GetUser(ctx, in.Uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by uuid: %v", err)
	}

	return &pkg.ResGetUserByUuid{
		Id: id,
	}, nil
}

func main() {
	ctx := context.Background()

	// MySQL DB 接続
	db, err := sql.Open("mysql", "docker:docker@tcp(localhost:3305)/test_database?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	// DB 接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	queries := tutorial.New(db)

	// サンプル: 全ユーザー取得
	users, err := queries.GetAllUsers(ctx)
	if err != nil {
		log.Printf("failed to get users: %v", err)
	} else {
		fmt.Printf("users: %+v\n", users)
	}

	// gRPC サーバー起動
	port := 8081
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// サービス登録。queries を持たせたサーバーを渡す
	pkg.RegisterCutListServiceServer(s, &server{queries: queries})

	log.Printf("start gRPC server port: %v", port)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
