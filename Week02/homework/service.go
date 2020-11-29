package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lunny/log"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func handleGetUserNameByID(c *gin.Context) {
	id := c.Param("id")
	idIntValue, err := strconv.Atoi(id)

	if err != nil {
		err = errors.Wrap(err, "invalid user id")
		log.Printf("%+v", err)
		// 参数类型错误，直接返回400
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	name, err := GetUserName(idIntValue)

	if err != nil {
		log.Printf("%+v", err)
		if errors.Is(err, sql.ErrNoRows) {
			// 如果是未找到此用户ID，直接返回404
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		// 其他错误降级处理
		c.String(http.StatusOK, "您好：游客")
		return
	}
	// 正常返回
	c.String(http.StatusOK, "您好：%s", name)
}
