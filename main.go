package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	pkg "kitta_backend/pkg"
	"kitta_backend/tutorial"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

type Cutlist struct {
	ID        int64
	ThisIsCut bool
	UserID    int64
	Name      string
	Color     string
	Count     int32
	Limit     int32
	LateTime  int32
	LateCount int32
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

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

func (s *server) CreateUser(ctx context.Context, in *pkg.ReqcreateUser) (*pkg.RescreateUser, error) {
	result, err := s.queries.CreateUser(ctx, in.Uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return &pkg.RescreateUser{
		Id: id,
	}, nil
}

func (s *server) CreateCutList(ctx context.Context, in *pkg.ReqCreateCutList) (*pkg.ResCreateCutList, error) {
	cutlistresult, err := s.queries.CreateCutList(ctx, tutorial.CreateCutListParams{
		ThisIsCut: in.ThisIsCut,
		UserID:    in.UserId,
		Name:      in.Name,
		Color:     in.Color,
		Count:     int32(in.Count),
		Limit:     int32(in.Limit),
		LateTime:  int32(in.LateTime),
		LateCount: int32(in.LateCount),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cut list: %v", err)
	}

	id, err := cutlistresult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %v", err)
	}

	historyresult, err := s.queries.CreateCutHistory(ctx, tutorial.CreateCutHistoryParams{
		ThisIsCut:      in.ThisIsCut,
		LateTime:       int32(in.LateTime),
		ListsID:        id,
		ListsUpdatedAt: time.Now(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create cut history: %v", err)
	}

	_ = historyresult // historyresult is not used further
	return &pkg.ResCreateCutList{
		Id: id,
	}, nil
}

// func (s *server) GetCutListByUserId(ctx context.Context, id int64) ([]Cutlist, error) {
// 	// This function is intentionally left unimplemented.
// 	cutlist, err := s.queries.GetCutLists(ctx, id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get cut lists by user id: %v", err)
// 	}

// 	return &pkg.ResGetCutList{
// 		ThisIsCut: cutlist.ThisIsCut,
// 		Id:        cutlist.UserID,
// 		Name:      cutlist.Name,
// 		Color:     cutlist.Color,
// 		Count:     int32(cutlist.Count),
// 		Limit:     int32(cutlist.Limit),
// 		LateTime:  int32(cutlist.LateTime),
// 		LateCount: int32(cutlist.LateCount),
// 	}

// 	// return cutlist, nil
// }

// func (s *server) GetCutList(ctx context.Context, userid int) (*pkg.ResGetCutListList, error) {
// 	// SQLC で取得
// 	cutlists, err := s.queries.GetCutLists(ctx, int64(userid))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get cut lists by user id: %v", err)
// 	}

// 	// gRPC の repeated 用スライス
// 	var pbCutLists []*pkg.ResGetCutList
// 	for _, c := range cutlists {
// 		pbCutLists = append(pbCutLists, &pkg.ResGetCutList{
// 			// state: 	nil,
// 			ThisIsCut: c.ThisIsCut,
// 			Id:        c.ID,
// 			Name:      c.Name,
// 			Color:     c.Color,
// 			Count:     int64(c.Count),
// 			Limit:     int64(c.Limit),
// 			// LateTime:  c.LateTime,
// 			LateCount: int64(c.LateCount),
// 			CreatedAt: c.CreatedAt.Time.Format(time.RFC3339),
// 			UpdatedAt: c.UpdatedAt.Time.Format(time.RFC3339),
// 			// unknownFields:    nil,
// 			// sizeCache:    0,
// 		})
// 	}

// 	return &pkg.ResGetCutListList{
// 		CutLists: pbCutLists,
// 	}, nil
// }

func (s *server) GetCutList(req *pkg.ReqGetCutList, stream pkg.CutListService_GetCutListServer) error {
	cutlists, err := s.queries.GetCutLists(stream.Context(), req.UserId)
	if err != nil {
		return fmt.Errorf("failed to get cut lists: %v", err)
	}

	for _, c := range cutlists {
		err := stream.Send(&pkg.ResGetCutList{
			ThisIsCut: c.ThisIsCut,
			Id:        c.ID,
			Name:      c.Name,
			Color:     c.Color,
			Count:     int64(c.Count),
			Limit:     int64(c.Limit),
			// LateTime:  int64(c.LateTime),
			LateCount: int64(c.LateCount),
			CreatedAt: c.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: c.UpdatedAt.Time.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}
	}

	return nil
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
