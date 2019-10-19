package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

const (
	defaultAvatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
)

type User struct {
	gorm.Model
	Name         string `gorm:"unique;not null"`
	Roles        string
	Token        string `gorm:"index:token"`
	Avatar       string
	Introduction string
	Email        string
}

type Routes struct {
	gorm.Model
	Name    string `gorm:"unique;not null"`
	Content string
}

func GetToken(username string) (token string, err error) {
	user := &User{Name: username}
	if err = DB.Where(user).First(user).Error; err != nil {
		return token, err
	}

	if user.Token == "" {
		uuids := strconv.Itoa(int(uuid.New().ID()))
		times := time.Now().UnixNano()
		token = fmt.Sprintf("T%d%v", times, uuids)
		err = DB.Model(user).Update("token", token).Error
		if err != nil {
			return token, err
		}

		return token, nil
	}

	return user.Token, nil
}

func CheckUser(token string) (user *User, err error) {
	user = &User{Token: token}
	if err = DB.Where(user).First(user).Error; err != nil {
		return user, err
	}

	return user, nil
}
