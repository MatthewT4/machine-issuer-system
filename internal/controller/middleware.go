package controller

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

	"machineIssuerSystem/pkg/jwt"
)

func (h *handlers) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := h.logger.With(
			slog.String("op", "AuthMiddleware"))

		cookie, err := c.Cookie(h.cfg.AuthCookieName)
		if err != nil {
			log.Warn("failed to get cookie: %v", err)

			return next(c)
		}

		jwtToken := cookie.Value

		claims, err := jwt.ParseToken(jwtToken, []byte(h.cfg.AuthSecretKey))
		if err != nil {
			log.Error("failed to parse token: %v", err)
		}

		role, ok := claims["role"].(float64)
		if !ok {
			log.Warn("no role in claims: %+v", claims)
		}

		id, ok := claims["id"].(string)
		if !ok {
			log.Warn("no id in claims: %+v", claims)
		}

		log.Info("setting role", role)
		log.Info("setting id", id)

		c.Set("role", role)
		c.Set("id", id)

		return next(c)
	}
}

func (h *handlers) PermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := h.logger.With(
			slog.String("op", "PermissionMiddleware"))

		method := c.Request().Method
		path := c.Request().URL.Path
		log.Info("params", method, path)

		if len(path) >= 5 && string(path[:5]) == "/auth" {
			return next(c)
		}

		ctxRole := c.Get("role")
		ctxID := c.Get("id")

		var role float64
		var id string

		if ctxRole != nil {
			role = ctxRole.(float64)
		}

		if ctxID != nil {
			id = ctxID.(string)
		}

		log.Info(fmt.Sprintf("user with role: %f, id: %s", role, id))

		//resp, err := h.core.GetPermissionHandler(c.Request().Context(), model.GetPermissionHandlerRequest{
		//	Method: method,
		//	Path:   path,
		//})
		//if err != nil {
		//	log.Error("failed to get permission: %v", err)
		//
		//	return err
		//}

		//if !lo.Contains(resp.Roles, int64(role)) {
		//	return echo.NewHTTPError(http.StatusForbidden, errorlist.ErrHandlerNotAllowed)
		//}

		return next(c)
	}
}
