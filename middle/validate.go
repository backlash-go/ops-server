package middle

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"ops-server/api"
	"ops-server/consts"
	"ops-server/db"
	"ops-server/utils"

    "ops-server/logs"
	"strings"
)

func BeforeRequestValidate(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		reqPath := c.Request().URL.Path

		if reqPath == "/api/ldap/user/auth" {
			return next(c)
		}

		token := c.Request().Header.Get("Authorization")

		userMap, err := db.RedisHGetAll(token)

		//if userMap == nil token不存在
		if userMap == nil {
			err := api.ErrorResp(c, consts.StatusText[consts.CodeNeedLogin], consts.CodeNeedLogin)
			return err
		}

		if err != nil {
			logs.GetLogger().Errorf("api BeforeRequestValidate RedisHGetAll is failed   err is %s\n", err.Error())
			err := api.ErrorResp(c, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
			return err
		}


		value, ok := userMap["roles"]
		if ok {
			for _, v := range strings.Split(value, ",") {
				if v == "admin" {
					return next(c)
				}
			}

			isPermission, err := utils.FilterPermission(strings.Split(value, ","), reqPath)
			fmt.Printf("isPermission is %v\n", isPermission)

			if err != nil {
				logs.GetLogger().Errorf("FilterPermission is faild err id %s\n", err)
				err := api.ErrorResp(c, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
				return err
			}
			if isPermission {
				return next(c)
			} else {
				err := api.ErrorResp(c, consts.StatusText[consts.CodeUserNoApiPermission], consts.CodeUserNoApiPermission)
				return err
			}
		} else {
			err := api.ErrorResp(c, consts.StatusText[consts.CodeUserNoAssignRole], consts.CodeUserNoAssignRole)
			return err
		}



		return next(c)

	}
}
