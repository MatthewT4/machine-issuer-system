package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"machineIssuerSystem/pkg/jwt"
)

func (h *handlers) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(h.cfg.Auth.CookieName)
		if err != nil {
			return err
		}

		jwtToken := cookie.Value

		claims, err := jwt.ParseToken(jwtToken, h.cfg.Auth.SecretKey)
		if err != nil {
			return err
		}

		role, ok := claims["role"].(int64)
		if !ok {
			h.logger.Warn("no role in cookie: %+v", claims)
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid role")
		}

		c.Set("role", role)

		return next(c)
	}
}
