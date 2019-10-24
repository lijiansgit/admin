package controllers

import (
	"fmt"
	"strings"

	"github.com/lijiansgit/admin/models"

	"github.com/gin-gonic/gin"
	"github.com/lijiansgit/admin/pkg/ldap"
)

var (
	User = &userController{}
)

type userController struct {
	Base
}

type LoginRequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseData struct {
	Token string `json:"token"`
}

func (u *userController) Login(c *gin.Context) {
	// username, password := c.PostForm("username"), c.PostForm("password")
	loginData := &LoginRequestPayload{}
	if err := c.BindJSON(loginData); err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	_, err := ldap.LDAP.Login(loginData.Username, loginData.Password)
	if err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	token, err := models.GetToken(loginData.Username)
	if err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	u.Base.composeJSON(c, &LoginResponseData{Token: token})
}

type UserResponseData struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles,omitempty"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
}

func (u *userController) Info(c *gin.Context) {
	user, err := models.CheckUser(c.Query("token"))
	if err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	users := &UserResponseData{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Roles:        u.Base.rolesToList(user.Roles),
		Introduction: user.Introduction,
		Avatar:       user.Avatar,
	}
	u.Base.composeJSON(c, users)
}

func (u *userController) Logout(c *gin.Context) {
	u.Base.composeJSON(c, "")
}

func (u *userController) List(c *gin.Context) {
	var (
		err       error
		users     []*models.User
		usersList []*UserResponseData
	)

	key := c.DefaultQuery("key", "0")
	// all users,uri: /list
	if key == "0" {
		users, err = models.GetAllUsers()
		if err != nil {
			u.Base.composeErrJSON(c, err)
			return
		}
	}

	// // search user by role, uri: /list?
	// if key != "0" {
	// 	id, err := strconv.Atoi(key)
	// 	if err != nil {
	// 		u.Base.composeErrJSON(c, err)
	// 		return
	// 	}
	// 	users, err = models.GetAllUsersByRoleID(id)
	// 	if err != nil {
	// 		u.Base.composeErrJSON(c, err)
	// 		return
	// 	}
	// }

	for _, user := range users {
		u := &UserResponseData{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			Roles:        u.Base.rolesToList(user.Roles),
			Introduction: user.Introduction,
			Avatar:       user.Avatar,
		}

		usersList = append(usersList, u)
	}

	u.Base.composeJSON(c, usersList)
}

type ModifyRolesRequestPayload struct {
	UserID uint     `json:"id"`
	Roles  []string `json:"roles,omitempty"`
}

func (u *userController) ModifyRoles(c *gin.Context) {
	var rolesData []*ModifyRolesRequestPayload
	if err := c.BindJSON(&rolesData); err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	for _, data := range rolesData {
		user := &models.User{}
		user.ID = data.UserID
		user.Roles = fmt.Sprintf("[%s]", strings.Join(data.Roles, ","))
		if err := models.ModifyRoles(user); err != nil {
			u.Base.composeErrJSON(c, err)
			return
		}
	}

	u.Base.composeJSON(c, "")
}

type SyncResponseData struct {
	Ok   int `json:"ok"`
	Fail int `json:"fail"`
}

func (u *userController) LDAPUserSyncDB(c *gin.Context) {
	users, err := ldap.LDAP.GetUsers()
	if err != nil {
		u.Base.composeErrJSON(c, err)
		return
	}

	resp := &SyncResponseData{}
	for name, email := range users {
		user := &models.User{
			Name:  name,
			Email: email,
		}

		if err := models.FirstOrCreate(user); err != nil {
			resp.Fail++
		} else {
			resp.Ok++
		}
	}

	u.Base.composeJSON(c, resp)
}
