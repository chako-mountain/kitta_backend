package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"reflect"

	_ "github.com/go-sql-driver/mysql"

	// "/tutorial"
	// "github.com/username/kitta_backend/tutorial"
	"kitta_backend/tutorial"

	"github.com/gin-gonic/gin"
)

func run() error {
	ctx := context.Background()

	r := gin.Default()

	// db, err := sql.Open("mysql", "docker:docker@/test_database?parseTime=true")
	db, err := sql.Open("mysql", "docker:docker@tcp(localhost:3305)/test_database?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	result, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}

	insertedAuthorID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Println(insertedAuthorID)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthorID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthorID, fetchedAuthor.ID))
	// return nil

	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r.Run()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

//テストです
