package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lunny/log"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func ServiceFactory() *Service {
	service := new(Service)
	service.dao = DAOFactory()
	return service
}

type Service struct {
	dao *DAO
}

func (s *Service) handleGetUserNameByID(c *gin.Context) {
	id := c.Param("id")
	idIntValue, err := strconv.Atoi(id)

	if err != nil {
		err = errors.Wrap(err, "invalid user id")
		log.Printf("%+v", err)
		// 参数类型错误，直接返回400
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	name, err := s.dao.GetUserName(idIntValue)

	if err != nil {
		log.Printf("%+v", err)
		if errors.Is(err, sql.ErrNoRows) {
			// 如果是未找到此用户ID，直接返回404
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			// 其他错误可降级处理
			c.String(http.StatusOK, "您好：%s", getDefaultUserName())
			// 或者返回500错误
			// c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	// 正常返回
	c.String(http.StatusOK, "您好：%s", name)
}

func getDefaultUserName() string {
	return "XXXX用户"
}
