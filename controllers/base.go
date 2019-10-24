package controllers

import (
	"strings"

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

func (b *Base) rolesToList(roles string) (list []string) {
	cutset := `[]'"`
	for _, v := range cutset {
		roles = strings.Replace(roles, string(v), "", -1)
	}

	list = strings.Split(roles, ",")
	if len(list) > 0 {
		return list
	}

	return nil
}
