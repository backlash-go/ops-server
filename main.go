package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"ops-server/api"
	"ops-server/db"
	"ops-server/middle"
)

func main() {
	e := echo.New()

	api.OperateLdap(e.Group("/api/ldap"))

	//设置跨域请求通行
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,echo.HeaderXRequestedWith,echo.HeaderAuthorization},
	}))

	e.Use(middle.BeforeRequestValidate)





	db.Init()
	db.InitLdap()

	e.Logger.Fatal(e.Start(":1323"))

}
