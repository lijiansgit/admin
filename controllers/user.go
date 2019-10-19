package controllers

import (
	"github.com/lijiansgit/admin/models"

	"github.com/gin-gonic/gin"
	"github.com/lijiansgit/admin/pkg/ldap"
)

var (
	User = &userController{}
)

type LoginResponseData struct {
	Token string `json:"token"`
}

type InfoResponseData struct {
	Name         string   `json:"name"`
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
}

type LoginRequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userController struct {
	Base
	LoginResponseData
	InfoResponseData
	LoginRequestPayload
}

func (u *userController) Login(c *gin.Context) {
	// username, password := c.PostForm("username"), c.PostForm("password")
	if err := c.BindJSON(&u.LoginRequestPayload); err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	_, err := ldap.LDAP.Login(u.Username, u.Password)
	if err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	u.LoginResponseData.Token, err = models.GetToken(u.Username)
	if err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	u.Base.composeJSON(c, u.LoginResponseData)
}

func (u *userController) Info(c *gin.Context) {
	user, err := models.CheckUser(c.Query("token"))
	if err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	u.InfoResponseData.Name = user.Name
	u.InfoResponseData.Roles = []string{"admin"}
	u.Base.composeJSON(c, u.InfoResponseData)
}

func (u *userController) Logout(c *gin.Context) {
	u.Base.composeJSON(c, u.InfoResponseData)
}

func (u *userController) Routes(c *gin.Context) {
}

func (u *userController) Roles(c *gin.Context) {
}
