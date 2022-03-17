package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ops-server/consts"
)

func SuccessResp(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"msg":  "ok",
		"code": consts.CodeSuccess,
		"data": data,
	})
}

func ErrorResp(ctx echo.Context, msg string, code int) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"msg":  msg,
		"code": code,
	})
}
