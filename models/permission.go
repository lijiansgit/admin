package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	Routes      string
}

func GetRoutes() (r string, err error) {
	role := &Role{}
	if err = DB.Where("name = ?", "system").First(role).Error; err != nil {
		return r, err
	}

	r = role.Routes
	return r, nil
}

func GetAllRoles() (roles []Role, err error) {
	err = DB.Find(&roles).Error
	if err != nil {
		return roles, err
	}

	return roles, nil
}

func GetAllRolesByName(name string) (roles []Role, err error) {
	err = DB.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).Find(&roles).Error
	if err != nil {
		return roles, err
	}

	return roles, nil
}

func GetAllRolesByID(id uint) (roles []Role, err error) {
	err = DB.First(&roles, id).Error
	if err != nil {
		return roles, err
	}

	return roles, nil
}

func CreateRole(role *Role) (key uint, err error) {
	if err = DB.Create(role).Error; err != nil {
		return key, err
	}

	if err = DB.First(role).Error; err != nil {
		return key, err
	}

	return role.ID, nil
}

func UpdateRole(role *Role) (err error) {
	if err = DB.Model(role).Updates(role).Error; err != nil {
		return err
	}

	return nil
}

func DeleteRole(role *Role) (err error) {
	// 软删除:https://gorm.io/zh_CN/docs/delete.html
	if err = DB.Delete(role).Error; err != nil {
		return err
	}

	return nil
}
