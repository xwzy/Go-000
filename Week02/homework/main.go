package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	service := ServiceFactory()
	router := gin.Default()
	router.GET("/user/:id", service.handleGetUserNameByID)
	router.Run(":8080")
}
