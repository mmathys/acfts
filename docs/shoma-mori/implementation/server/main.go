package main

import (
	"acfts/api"
	"acfts/db/model"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// n is the number of servers
const n = 4

func initDB(port int) *gorm.DB {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/acfts_"+strconv.Itoa(port)+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Output{})

	return db
}

func initRoute(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// APIs
	r.GET("/", api.Ping())
	r.POST("/genesis", api.CreateGenesis(db))
	r.POST("/transaction", api.VerifyTransaction(db, n))
	r.DELETE("/all", api.DeleteAll(db))

	return r
}

func main() {
	fmt.Printf("Input port number: ")
	var port int
	fmt.Scan(&port)

	db := initDB(port)
	defer db.Close()

	go func() {
		log.Println(http.ListenAndServe("localhost:8000", nil))
	}()

	r := initRoute(db)
	r.Run(":" + strconv.Itoa(port))
}
