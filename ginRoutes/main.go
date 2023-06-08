package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Photo struct {
	AlbumId      int    `json:"albumId"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

func getAllPhotos(c *gin.Context) {

	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Cannot connect to database")
		log.Println(err)
	}

	var photoList []Photo

	queryRes, err := db.Raw("select * from photos").Rows()
	if err != nil {
		log.Println(err)
	}

	defer queryRes.Close()

	for queryRes.Next() {
		var p Photo
		db.ScanRows(queryRes, &p)
		photoList = append(photoList, p)
	}

	// fmt.Println(photoList)

	c.IndentedJSON(http.StatusOK, photoList)
}

func getUserPhotos(c *gin.Context) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	id := c.Param("id")
	queryRes, err := db.Raw("select * from photos where album_id = ?", id).Rows()
	if err != nil {
		log.Panicln(err)
	}

	var userP []Photo
	defer queryRes.Close()

	for queryRes.Next() {
		var p Photo
		db.ScanRows(queryRes, &p)
		userP = append(userP, p)
	}

	c.IndentedJSON(http.StatusOK, userP)
}

func main() {
	fmt.Println("Making http endpoints using Gin")
	router := gin.Default()
	router.GET("/photos", getAllPhotos)
	router.GET("/photos/:id", getUserPhotos)

	log.Fatal(router.Run("localhost:8081"), nil)
}
