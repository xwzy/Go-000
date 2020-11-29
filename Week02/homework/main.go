package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/user/:id", handleGetUserNameByID)
	router.Run(":8080")
}
