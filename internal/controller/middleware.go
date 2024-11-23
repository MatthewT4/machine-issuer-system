package controller

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"machineIssuerSystem/internal/model"
	"machineIssuerSystem/pkg/errorlist"
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

		err = next(c)
		resp := c.Response()
		resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		resp.Header().Set("Access-Control-Allow-Headers", "*")
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
		c.SetResponse(resp)
		return err
	}
}

func (h *handlers) PermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := h.logger.With(
			slog.String("op", "PermissionMiddleware"))

		method := c.Request().Method
		path := c.Request().URL.Path
		log.Info("params", method, path)

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

		resp, err := h.core.GetPermissionHandler(c.Request().Context(), model.GetPermissionHandlerRequest{
			Method: method,
			Path:   path,
		})
		if err != nil {
			log.Error("failed to get permission: %v", err)

			return err
		}

		if !lo.Contains(resp.Roles, int64(role)) {
			return echo.NewHTTPError(http.StatusForbidden, errorlist.ErrHandlerNotAllowed)
		}

		return next(c)
	}
}
