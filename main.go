package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	// "log"

	// "reflect"
	"context"

	_ "github.com/go-sql-driver/mysql"

	// "/tutorial"
	// "github.com/username/kitta_backend/tutorial"
	pkg "kitta_backend/pkg"
	"kitta_backend/tutorial"

	// "github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	pkg.UnimplementedCutListServiceServer
}

// func (s *server) GetCutLists(ctx context.Context, userID int64) (*pkg.GetCutListsResponse, error) {
// 	// Implement your logic to retrieve cut lists here
// 	return &pkg.GetCutListsResponse{}, nil
// }

func (s *server) GetUserByUuid(ctx context.Context, in *pkg.ReqGetUserByUuid) (*pkg.ResGetUserByUuid, error) {
	return &pkg.ResGetUserByUuid{
		Id: 1, // ← proto の定義通り
	}, nil
}

// func (s *server) GetUserByUuid(ctx context.Context, in *pkg.ReqGetUserByUuid) (*pkg.ResGetUserByUuid, error) {
// 	queries := tutorial.New(db) // ここはサーバに DB を持たせるとよい
// 	user, err := queries.GetUserByUuid(ctx, in.Uuid)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get user by uuid: %v", err)
// 	}

// 	return &pkg.ResGetUserByUuid{
// 		Id: user.ID,
// 	}, nil
// }

func main() {
	// userID := int64(1)

	ctx := context.Background()

	// r := gin.Default()

	// db, err := sql.Open("mysql", "docker:docker@/test_database?parseTime=true")
	db, err := sql.Open("mysql", "docker:docker@tcp(localhost:3305)/test_database?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		fmt.Print(err)
	}

	port := 8081
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	queries := tutorial.New(db)
	// newUUID := uuid.New()
	// uuidStr := newUUID.String()
	// err = queries.CreateUser(ctx, newUUID.String())
	// if err != nil {
	// 	fmt.Print(err)
	// }

	fmt.Print(queries.GetAllUsers(ctx))

	s := grpc.NewServer()

	// return nil

	log.Printf("start gRPC server port: %v", port)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// r.GET("/ping", func(c *gin.Context) {
	// 	// Return JSON response
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// r.Run()
	pkg.RegisterCutListServiceServer(s, &server{})

	log.Printf("start gRPC server port: %v", port)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}

}

//テストです

// package main

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net"

// 	pkg "kitta_backend/pkg"
// 	"kitta_backend/tutorial"

// 	_ "github.com/go-sql-driver/mysql"
// 	"google.golang.org/grpc"
// )

// type server struct {
// 	pkg.UnimplementedCutListServiceServer
// }

// // gRPC メソッドのサンプル実装
// func (s *server) GetUserByUuid(ctx context.Context, in *pkg.ReqGetUserByUuid) (*pkg.ResGetUserByUuid, error) {
// 	return &pkg.ResGetUserByUuid{
// 		Id: 1, // proto の定義に合わせたサンプル
// 	}, nil
// }

// func main() {
// 	ctx := context.Background()

// 	// MySQL DB 接続
// 	db, err := sql.Open("mysql", "docker:docker@tcp(localhost:3305)/test_database?parseTime=true&loc=Asia%2FTokyo")
// 	if err != nil {
// 		log.Fatalf("failed to connect to DB: %v", err)
// 	}
// 	defer db.Close()

// 	// DB 接続確認
// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("DB ping failed: %v", err)
// 	}

// 	queries := tutorial.New(db)

// 	// サンプル: 全ユーザー取得
// 	users, err := queries.GetAllUsers(ctx)
// 	if err != nil {
// 		log.Printf("failed to get users: %v", err)
// 	} else {
// 		fmt.Printf("users: %+v\n", users)
// 	}

// 	// gRPC サーバー起動
// 	port := 8081
// 	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()

// 	// サービス登録
// 	pkg.RegisterCutListServiceServer(s, &server{})

// 	log.Printf("start gRPC server port: %v", port)
// 	if err := s.Serve(listener); err != nil {
// 		log.Fatalf("failed to serve gRPC: %v", err)
// 	}
// }
