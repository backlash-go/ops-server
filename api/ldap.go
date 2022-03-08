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
	g.POST("/user/create", CreateUser)

}



func CreateUser(c echo.Context) error {
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

	if req.Cn == "" || req.DisplayName == "" || req.GivenName == "" || req.EmployeeType == "" || req.UserPassword == "" {
		return c.String(http.StatusBadRequest, `req is null maybe   : req.cn == "" || req.displayName  == ""|| req.givenName  == "" || req.employeeType == ""`)
	}

	if err := DefaultLdap.CreateLdapUser(req); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("CreateLdapUser request is failed err is :%s\n", err))
	}

	return c.String(http.StatusOK, "CreateUser  is successful")

}
