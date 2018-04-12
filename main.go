package main

import (
	"context"
	"flag"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smallnest/rpcx/client"
)

var (
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_users/Users", "prefix path")
	xclient  client.XClient
)

type Users struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	PID    int
	Name   string
	Gender int
	Mark   string
}

func main() {
	// init RPC
	d := client.NewEtcdDiscovery(*basePath, "", []string{*etcdAddr}, nil)
	xclient = client.NewXClient("Users", client.Failover, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	// init gin
	r := gin.Default()

	// Get all users
	r.GET("/users", func(c *gin.Context) {
		users := &[]Users{}
		err := xclient.Call(context.Background(), "GetAllUsers", &Users{}, users)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":     err.Error(),
				"success": false,
			})
		} else {
			c.JSON(200, gin.H{
				"data":    users,
				"msg":     "OK",
				"success": true,
			})
		}
	})

	// Get user by id
	r.GET("/users/:id", func(c *gin.Context) {
		u64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		args := &Users{ID: u64}
		result := &Users{}

		// err := DB.Get(&user, "SELECT id,name,gender,mark,created_at FROM users WHERE id=?", id)
		err := xclient.Call(context.Background(), "GetUser", args, result)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":     err.Error(),
				"success": false,
			})
		} else {
			c.JSON(200, gin.H{
				"data":    result,
				"msg":     "OK",
				"success": true,
			})
		}
	})

	// Create user
	r.POST("/users", func(c *gin.Context) {
		user := &Users{}
		c.ShouldBindWith(user, binding.FormPost)
		err := xclient.Call(context.Background(), "AddUser", user, nil)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":     err.Error(),
				"success": false,
			})
		} else {
			c.JSON(200, gin.H{
				"msg":     "OK",
				"success": true,
			})
		}
	})

	// Update user
	r.PUT("/users", func(c *gin.Context) {
		msg := gin.H{"msg": "OK", "success": true}
		user := &Users{}
		c.ShouldBindWith(user, binding.FormPost)
		if user.ID != 0 {
			err := xclient.Call(context.Background(), "UpdateUser", user, nil)
			if err != nil {
				msg["success"] = false
				msg["msg"] = err.Error()
			}
		} else {
			msg["success"] = false
			msg["msg"] = "ID Can not Empty"
		}
		c.JSON(200, msg)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
