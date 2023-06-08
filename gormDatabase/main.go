package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//https://jsonplaceholder.typicode.com/photos

type Photos struct {
	AlbumId      int    `json:"albumId"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

func main() {
	fmt.Println("API for Photos")

	resp, err := http.Get("https://jsonplaceholder.typicode.com/photos")
	if err != nil {
		log.Println(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(string(content))

	var photos []Photos
	if err := json.Unmarshal(content, &photos); err != nil {
		log.Println("cannot Unmarshall ", err)
	}

	for _, p := range photos {
		fmt.Println("The Album id is: ", p.AlbumId)
		fmt.Println("The Id is: ", p.Id)
		fmt.Println("The Title is: ", p.Title)
		fmt.Println("The Url is: ", p.Url)
		fmt.Println("The ThumbnailUrl is: ", p.ThumbnailUrl)
		fmt.Println("")
	}

	fmt.Println("The Unmarshalling and storing data is done")

	//connecting to the database using gorm
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Cannot connect to database")
		log.Println(err)
	}

	// fmt.Println(db)
	if err := db.Exec("create table photos (album_id int, id int, title varchar(255), url varchar(500), thumbnail_url varchar(500))"); err != nil {
		log.Println(err)
	}

	fmt.Println("Table created")

	for _, p := range photos {
		db.Exec("insert into photos(album_id, id, title, url, thumbnail_url) values (?,?,?,?,?)", p.AlbumId, p.Id, p.Title, p.Url, p.ThumbnailUrl)
	}

	fmt.Println("Data inserted")

}
