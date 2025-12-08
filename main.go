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
	"kitta_backend/tutorial"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

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
	newUUID := uuid.New()
	// uuidStr := newUUID.String()
	err = queries.CreateUser(ctx, newUUID.String())
	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(queries.GetAllUsers(ctx))

	s := grpc.NewServer()

	// return nil

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// r.GET("/ping", func(c *gin.Context) {
	// 	// Return JSON response
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// r.Run()

}

//テストです
