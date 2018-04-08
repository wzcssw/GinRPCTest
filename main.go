package main

import (
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

type Users struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Gender    int       `db:"gender"`
	Mark      string    `db:"mark"`
	CreatedAt time.Time `db:"created_at"`
}

func init() {
	DB, _ = sqlx.Connect("mysql", "root:wzc19931030@tcp(127.0.0.1:3306)/funny?charset=utf8&parseTime=true")
}

func main() {
	r := gin.Default()

	r.GET("/users", func(c *gin.Context) {
		users := []Users{}
		DB.Select(&users, "SELECT id,name,gender,mark,created_at FROM users")
		c.JSON(200, users)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		user := Users{}
		id := c.Param("id")
		err := DB.Get(&user, "SELECT id,name,gender,mark,created_at FROM users WHERE id=?", id)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":     err.Error(),
				"success": false,
			})
		} else {
			c.JSON(200, gin.H{
				"data":    user,
				"msg":     "OK",
				"success": true,
			})
		}

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
