package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"machineIssuerSystem/internal/model"
	"machineIssuerSystem/pkg/errorlist"
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
			return echo.NewHTTPError(http.StatusUnauthorized, errorlist.ErrInvalidRole)
		}

		c.Set("role", role)

		return next(c)
	}
}

func (h *handlers) PermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method
		path := c.Request().URL.Path
		role := c.Get("role").(int64)

		resp, err := h.core.GetPermissionHandler(c.Request().Context(), model.GetPermissionHandlerRequest{
			Method: method,
			Path:   path,
		})
		if err != nil {
			return err
		}

		// handler is public if no roles added
		if len(resp.Roles) == 0 {
			return next(c)
		}

		if !lo.Contains(resp.Roles, role) {
			return echo.NewHTTPError(http.StatusForbidden, errorlist.ErrHandlerNotAllowed)
		}

		return next(c)
	}
}
