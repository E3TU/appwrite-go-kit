// Command gin is a runnable example of appwrite-base.
//
// It shows the whole chain in one file: build a Config in code, call base.New,
// then guard a route with the auth middleware and read a row through the
// database service.
package main

import (
	"log"
	"net/http"

	"github.com/appwrite/sdk-for-go/v7/query"
	"github.com/e3tu/appwrite-base"
	"github.com/e3tu/appwrite-base/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := base.Config{
		ProjectID:      "example-project",
		APIKey:         "standard_...",
		Endpoint:       "http://appwrite.local/v1",
		PublicEndpoint: "https://appwrite.example.com/v1",
		Auth:           &base.AuthConfig{},
		DB:             &base.DBConfig{DatabaseID: "main"},
	}

	b, err := base.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	if a, ok := b.Auth(); ok {
		r.GET("/me", a.Middleware(""), func(c *gin.Context) {
			u, _ := auth.UserFrom(c)
			c.JSON(http.StatusOK, u)
		})
	}

	if db, ok := b.DB(); ok {
		r.GET("/notes", func(c *gin.Context) {
			rows, err := db.TablesDB().ListRows(
				db.DatabaseID(),
				"notes",
				db.TablesDB().WithListRowsQueries([]string{
					query.OrderDesc("$createdAt"),
					query.Limit(25),
				}),
			)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, rows)
		})
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
