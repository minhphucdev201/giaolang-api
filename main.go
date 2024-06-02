package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Room struct {
	ID           int    `json:"id"`
	Code         string `json:"code"`
	Row_of_house string `json:"row_of_house"`
}

var rooms []Room

func main() {
	db, err := sql.Open("mysql", "root:mambau2001@tcp(127.0.0.1:3306)/defaultdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected to database", db)
	r := gin.Default()
	r.GET("/rooms", func(c *gin.Context) {
		getAll(c, db)
	})
	r.POST("/room", func(c *gin.Context) {
		create(c, db)
	})
	// Other endpoints...

	r.Run()
}

func create(c *gin.Context, db *sql.DB) {
	var newRoom Room
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO rooms (id, code, row_of_house) VALUES (  ?, ?, ?)",
		newRoom.ID, newRoom.Code, newRoom.Row_of_house)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	newRoom.ID = int(id)
	c.JSON(201, newRoom)
}

func getAll(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM rooms")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {

		var room Room
		err := rows.Scan(
			&room.ID, &room.Code, &room.Row_of_house,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, rooms)
}

func generateID() string {
	// Sử dụng thời gian hiện tại để tạo ID ngẫu nhiên
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
