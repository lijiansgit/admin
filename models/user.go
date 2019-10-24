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

func FirstOrCreate(user *User) (err error) {
	// todo 有更新用户邮箱等信息的情况
	if err = DB.Create(user).Error; err != nil {
		return err
	}

	return nil
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
	if err = DB.Where("token = ?", token).First(user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetAllUsers() (users []*User, err error) {
	if err = DB.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

// func GetAllUsersByRoleID(roleID int) (users []*User, err error) {
// 	err = DB.Where("roles LIKE ?", fmt.Sprintf("%%%d%%", roleID)).Find(&users).Error
// 	if err != nil {
// 		return users, err
// 	}

// 	return users, nil
// }

func ModifyRoles(user *User) (err error) {
	if err = DB.Model(user).Updates(user).Error; err != nil {
		return err
	}

	return nil
}
