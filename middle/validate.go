package middle

import "github.com/labstack/echo/v4"

func BeforeRequestValidate(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		authValue := c.Request().Header.Get("Authorization")
		if authValue == "677cbf56dfa43b8854" {
			return next(c)
		}
		return c.JSON(403, "需要授权验证")
	}
}
