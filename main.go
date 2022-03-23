package main

import (
	"github.com/labstack/echo/v4"
	"ops-server/api"
	"ops-server/db"
	"ops-server/middle"
)

func main() {
	e := echo.New()

	api.OperateLdap(e.Group("/api/ldap"))

	e.Use(middle.BeforeRequestValidate)


	db.Init()
	db.InitLdap()

	e.Logger.Fatal(e.Start(":1323"))

}
