package main

import (
	"backend/database/dbInfo"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// var dataResult []result

const (
	dbUserName = "kevinwu"
	dbPassword = "xxxxxxx"
	dbIp       = "xxxxxxxx"
	dbName     = "FAFT_TEST"
	tableName  = "Result"
)

type result struct {
	Id         int    `json:"id"`
	Time       string `json:"time"`
	Tester     string `json:"tester"`
	Name       string `json:"name"`
	Board      string `json:"board"`
	Model      string `json:"model"`
	Version    string `json:"version"`
	LogPath    string `json:"logPath"`
	PassOrFail string `json:"passOrFail"`
	Reason     string `json:"reason"`
}

var dataResult []dbInfo.Result

func main() {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUserName, dbPassword, dbIp, dbName))
	if err != nil {
		log.Fatal("Failed to open db: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	log.Print("Database is connected")

	var err2 error
	router.GET("/test", func(ctx *gin.Context) {
		dataResult, err2 = dbInfo.SearhData(db)
		if err2 != nil {
			log.Fatal("Failed to search data in table `Result`: ", err)
		}
		ctx.IndentedJSON(http.StatusOK, dataResult)
		dataResult = nil
	})

	router.Static("/logDB", "/home/ubuntu/backend/logDB")
	router.Run(":8082")
}

// Look in the future
// func DownloadLicense(ctx *gin.Context) {
// 	content := "Download file here happliy"
// 	fileName := "/log/log.txt"
// 	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
// 	ctx.Header("Content-Type", "application/text/plain")
// 	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
// 	ctx.Writer.Write([]byte(content))
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"msg": "Download file successfully",
// 	})
// }
