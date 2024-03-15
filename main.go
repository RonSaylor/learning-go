package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/go-sql-driver/mysql"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("mysqluser"),
		Passwd: os.Getenv("mysqlpass"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "albums",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.PUT("/albums", updateAlbum)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	var albums []album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no albums found"})
	}
	defer rows.Close()

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no albums found"})
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no albums found"})
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	var alb album

	rows, err := db.Query("SELECT * FROM album WHERE id = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no albums found"})
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no albums found"})
			return
		}
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no albums found"})
	}

	c.IndentedJSON(http.StatusOK, alb)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})

	}

	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, id)
}

func updateAlbum(c *gin.Context) {
	var albumToUpdate album

	if err := c.BindJSON(&albumToUpdate); err != nil {
		return
	}

	_, err := db.Exec("UPDATE album SET title = ?, artist = ?, price = ? WHERE id = ?", albumToUpdate.Title, albumToUpdate.Artist, albumToUpdate.Price, albumToUpdate.ID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, albumToUpdate)
}
