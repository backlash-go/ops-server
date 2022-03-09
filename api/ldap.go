package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"ops-server/entity"
	"ops-server/service"
)

func OperateLdap(g *echo.Group) {
	g.POST("/user/create", CreateLdapUser)
	g.DELETE("/user/delete", DeleteLdapUser)

}

func CreateLdapUser(c echo.Context) error {
	//create ldap connection
	DefaultLdap, err := service.CreateLdapConnection()
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("CreateLdapConnection is failed err is :%s\n", err))
	}

	// receive create user http request params
	req := &entity.CreateUserParams{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("c.Bind(req) is failed err is %s\n", err))
	}

	log.Printf("req is  %s\n", req)

	if req.Cn == "" || req.DisplayName == "" || req.GivenName == "" || len(req.EmployeeType) == 0 || req.UserPassword == "" {
		return c.String(http.StatusBadRequest, `req is null maybe   : req.cn == "" || req.displayName  == ""|| req.givenName  == "" || req.employeeType == 0`)
	}

	if err := DefaultLdap.CreateUser(req); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("CreateLdapUser request is failed err is :%s\n", err))
	}

	return c.String(http.StatusOK, "CreateUser  is successful")

}

func DeleteLdapUser(c echo.Context) error {
	//create ldap connection
	DefaultLdap, err := service.CreateLdapConnection()
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("CreateLdapConnection is failed err is :%s\n", err))
	}

	req := new(entity.DeleteUserParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("c.Bind(req) is failed err is %s\n", err))
	}

	if req.Dn == "" {
		return c.String(http.StatusBadRequest, `Delete ldap user req is null maybe   : req.cn == ""`)
	}

	if err := DefaultLdap.DeleteUser(req); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("DeleteUser request is failed err is :%s\n", err))
	}

	return c.String(http.StatusOK, "DeleteUser  is successful")

}
