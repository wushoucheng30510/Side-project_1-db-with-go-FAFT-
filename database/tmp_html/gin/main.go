package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUserName = "root"
	dbPassword = "A19930107"
	dbIp       = "localhost"
	dbName     = "test"
)

type result struct{
	id int 
	name string
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("temp/*.html")
	router.GET("/",func(ctx *gin.Context){
		ctx.HTML(200, "index.html", gin.H{})
	})

	router.GET("/data",func(ctx *gin.Context){
		ctx.HTML(200, "data.html", gin.H{"data": "Hello Go/Gin world."})
	})

	router.GET("/form",func(ctx *gin.Context){
		ctx.HTML(200, "form.html", gin.H{})
	})

	router.POST("/service",func(ctx *gin.Context){
		testName := ctx.PostForm("testName")
		ctx.JSON(200,gin.H{
			"result": "ok",
			"hello":testName,
		})
	})

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUserName, dbPassword, dbIp, dbName))
	if err != nil {
		log.Fatal("Failed to open db: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	log.Print("Database is connected")


	SearhData(db,"Kevin")

	router.Run(":8080")
}

func SearhData(db *sql.DB, name string){
	result := new(result)
	row := db.QueryRow("select * from result where name=?", name)
	if err := row.Scan(&result.id,&result.name);err!=nil{
		log.Print("Failed to search data: ",err)
		return 
	}
	log.Print("Data: ",*result)
}
