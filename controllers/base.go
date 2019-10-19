package controllers

import (
	"github.com/gin-gonic/gin"
	log "github.com/lijiansgit/go/libs/log4go"
)

type Base struct {
}

func (b *Base) composeErrJSON(c *gin.Context, err error) {
	errMsg := err.Error()
	log.Warn("%s: %s", c.Request.RequestURI, errMsg)
	c.JSON(200, map[string]interface{}{
		"code": errCode,
		"msg":  errMsg,
		"data": "",
	})
}

func (b *Base) composeJSON(c *gin.Context, data interface{}) {
	c.JSON(200, map[string]interface{}{
		"code": okCode,
		"msg":  "",
		"data": data,
	})
}
