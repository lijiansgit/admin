package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lijiansgit/admin/models"
)

var (
	Permission = &permissionController{}
)

type Route struct {
	Path       string     `json:"path"`
	Component  string     `json:"component"`
	Hidden     bool       `json:"hidden"`
	Redirect   string     `json:"redirect"`
	AlwaysShow bool       `json:"alwaysShow"`
	Children   []Children `json:"children,omitempty"`
	Meta       Meta       `json:"meta"`
}

type Children struct {
	Path      string     `json:"path"`
	Component string     `json:"component"`
	Name      string     `json:"name"`
	Redirect  string     `json:"redirect"`
	Children  []Children `json:"children,omitempty"`
	Meta      Meta       `json:"meta"`
}

type Meta struct {
	Title   string   `json:"title"`
	NoCache bool     `json:"noCache"`
	Icon    string   `json:"icon"`
	Roles   []string `json:"roles,omitempty"`
}

type Role struct {
	Key         uint    `json:"key"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Routes      []Route `json:"routes,omitempty"`
}

type permissionController struct {
	Base
	RoleRequestPayload
	RoleResponseData
}

func (p *permissionController) Routes(c *gin.Context) {
	routesStr, err := models.GetRoutes()
	if err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	var routes []Route
	err = json.Unmarshal([]byte(routesStr), &routes)
	if err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	p.Base.composeJSON(c, routes)
}

func (p *permissionController) Roles(c *gin.Context) {
	var (
		err   error
		r     []models.Role
		roles []*Role
	)

	name, keyStr := c.DefaultQuery("name", "0"), c.DefaultQuery("key", "0")
	// all roles, uri: /roles
	if name == keyStr {
		r, err = models.GetAllRoles()
		if err != nil {
			p.Base.composeErrJSON(c, err)
			return
		}
	}

	// role search by name, uri: /roles?name=name
	if name != "0" {
		r, err = models.GetAllRolesByName(name)
		if err != nil {
			p.Base.composeErrJSON(c, err)
			return
		}
	}

	// role search by id, uri: /roles?key=key
	if keyStr != "0" {
		key, err := strconv.Atoi(keyStr)
		if err != nil {
			p.Base.composeErrJSON(c, err)
			return
		}

		r, err = models.GetAllRolesByID(uint(key))
		if err != nil {
			p.Base.composeErrJSON(c, err)
			return
		}
	}

	for _, v := range r {
		role := new(Role)
		role.Key = v.ID
		role.Name = v.Name
		role.Description = v.Description
		err := json.Unmarshal([]byte(v.Routes), &role.Routes)
		if err != nil {
			p.Base.composeErrJSON(c, err)
			return
		}

		roles = append(roles, role)
	}

	p.Base.composeJSON(c, roles)
}

type RoleRequestPayload struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Routes      []Route `json:"routes,omitempty"`
}

type RoleResponseData struct {
	Key uint `json:"key"`
}

func (p *permissionController) CreateRole(c *gin.Context) {
	if err := c.BindJSON(&p.RoleRequestPayload); err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	routesByte, err := json.Marshal(p.RoleRequestPayload.Routes)
	if err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	role := &models.Role{
		Name:        p.RoleRequestPayload.Name,
		Description: p.RoleRequestPayload.Description,
		Routes:      string(routesByte),
	}

	p.RoleResponseData.Key, err = models.CreateRole(role)
	if err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	p.Base.composeJSON(c, p.RoleResponseData)
}

func (p *permissionController) UpdateRole(c *gin.Context) {
	keyStr := c.Param("key")
	key, err := strconv.Atoi(keyStr)
	if err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	if err := c.BindJSON(&p.RoleRequestPayload); err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	b, err := json.Marshal(p.RoleRequestPayload.Routes)
	if err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	role := &models.Role{
		Name:        p.RoleRequestPayload.Name,
		Description: p.RoleRequestPayload.Description,
		Routes:      string(b),
	}
	role.ID = uint(key)
	if err := models.UpdateRole(role); err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	p.Base.composeJSON(c, "")
}

func (p *permissionController) DeleteRole(c *gin.Context) {
	keyStr := c.Param("key")
	key, err := strconv.Atoi(keyStr)
	if err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	role := &models.Role{}
	role.ID = uint(key)
	if err := models.DeleteRole(role); err != nil {
		p.Base.composeErrJSON(c, err)
		return
	}

	p.Base.composeJSON(c, "")
}
