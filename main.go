package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/ldap.v2"
	"ops-server/api"


)

func main() {
	e := echo.New()


	api.OperateLdap(e.Group("/ldap"))

	e.Logger.Fatal(e.Start(":1323"))

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "ldap.example.com", 389))

	l.Bind()




}
