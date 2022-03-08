package main

import (
	"github.com/labstack/echo/v4"
	"ops-server/api"
	"ops-server/middle"
)

func main() {
	e := echo.New()

	api.OperateLdap(e.Group("/ldap"))

	e.Use(middle.BeforeRequestValidate)

	e.Logger.Fatal(e.Start(":1323"))





}
