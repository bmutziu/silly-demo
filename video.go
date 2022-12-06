package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-pg/pg/extra/pgotel/v10"
	"github.com/go-pg/pg/v10"

	"github.com/gin-gonic/gin"
)

// Test:
// $ docker container run --name my-db -e POSTGRES_PASSWORD=postgres -d --publish 5432:5432 postgres
// $ DB_ENDPOINT=127.0.0.1 DB_PORT=5432 DB_USER=postgres DB_PASS=postgres DB_NAME=postgres go run .
// $ curl -X POST "http://localhost:8080/video?id=wNBG1-PSYmE&title=Kubernetes%20Policies%20And%20Governance%20-%20Ask%20Me%20Anything%20With%20Jim%20Bugwadia"
// $ curl -X POST "http://localhost:8080/video?id=VlBiLFaSi7Y&title=Scaleway%20-%20Everything%20We%20Expect%20From%20A%20Cloud%20Computing%20Service%3F"
// $ curl "http://localhost:8080/videos" | jq .

var dbSession *pg.DB = nil

type Video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func videosGetHandler(c *gin.Context) {
	traceContext, _ := tp.Tracer(name).Start(context.TODO(), name)

	db := getDB(c)
	if db == nil {
		return
	}
	var videos []Video
	err := db.ModelContext(traceContext, &videos).Select()
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, videos)
}

func videoPostHandler(c *gin.Context) {
	traceContext, _ := tp.Tracer(name).Start(context.TODO(), "silly-demo4")

	db := getDB(c)
	if db == nil {
		return
	}
	id := c.Query("id")
	if len(id) == 0 {
		fmt.Println("id is empty")
		c.String(http.StatusBadRequest, "id is empty")
		return
	}
	title := c.Query("title")
	if len(title) == 0 {
		fmt.Println("title is empty")
		c.String(http.StatusBadRequest, "title is empty")
		return
	}
	video := &Video{
		ID:    id,
		Title: title,
	}
	_, err := db.ModelContext(traceContext, video).Insert()
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}

func getDB(c *gin.Context) *pg.DB {
	if dbSession != nil {
		return dbSession
	}
	endpoint := os.Getenv("DB_ENDPOINT")
	if len(endpoint) == 0 {
		fmt.Println("Environment variable `DB_ENDPOINT` is empty")
		c.String(http.StatusBadRequest, "Environment variable `DB_ENDPOINT` is empty")
		return nil
	}
	port := os.Getenv("DB_PORT")
	if len(port) == 0 {
		fmt.Println("Environment variable `DB_PORT` is empty")
		c.String(http.StatusBadRequest, "Environment variable `DB_PORT` is empty")
		return nil
	}
	user := os.Getenv("DB_USER")
	if len(user) == 0 {
		user = os.Getenv("DB_USERNAME")
		if len(user) == 0 {
			fmt.Println("Environment variables `DB_USER` and `DB_USERNAME` are empty")
			c.String(http.StatusBadRequest, "Environment variables `DB_USER` and `DB_USERNAME` are empty")
			return nil
		}
	}
	pass := os.Getenv("DB_PASS")
	if len(pass) == 0 {
		pass = os.Getenv("DB_PASSWORD")
		if len(pass) == 0 {
			fmt.Println("Environment variables `DB_PASS` and `DB_PASSWORD are empty")
			c.String(http.StatusBadRequest, "Environment variables `DB_PASS` and `DB_PASSWORD are empty")
			return nil
		}
	}
	name := os.Getenv("DB_NAME")
	if len(name) == 0 {
		fmt.Println("Environment variable `DB_NAME` is empty")
		c.String(http.StatusBadRequest, "Environment variable `DB_NAME` is empty")
		return nil
	}
	dbSession := pg.Connect(&pg.Options{
		Addr:     endpoint + ":" + port,
		User:     user,
		Password: pass,
		Database: name,
	})
	dbSession.AddQueryHook(pgotel.NewTracingHook())
	return dbSession
}
