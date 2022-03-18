package middle

import (
	"github.com/labstack/echo/v4"
	"ops-server/api"
	"ops-server/consts"
	"ops-server/db"
)

func BeforeRequestValidate(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		token := c.Request().Header.Get("Authorization")

		_, err := db.RedisHMGet(token, "id", "user_name", "email")

		if err != nil {
			api.ErrorResp(c, consts.StatusText[consts.CodeNeedLogin], consts.CodeNeedLogin)
		}

		return next(c)

	}
}
