package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func OperateLdap(g *echo.Group) {
	g.GET("/user/add", CreateUser)

}

func CreateUser(c echo.Context) error {

	return c.String(http.StatusOK, "sad11")

}
