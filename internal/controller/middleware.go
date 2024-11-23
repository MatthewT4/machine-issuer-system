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

		log.Info("setting role", role)

		c.Set("role", role)

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
		ctxRole := c.Get("role")

		var role float64
		if ctxRole != nil {
			role = ctxRole.(float64)
		}

		log.Info(fmt.Sprintf("user with role: %d", role))

		resp, err := h.core.GetPermissionHandler(c.Request().Context(), model.GetPermissionHandlerRequest{
			Method: method,
			Path:   path,
		})
		if err != nil {
			log.Error("failed to get permission: %v", err)

			return err
		}

		// handler is public if no roles added
		if len(resp.Roles) == 0 {
			return next(c)
		}

		if !lo.Contains(resp.Roles, int64(role)) {
			return echo.NewHTTPError(http.StatusForbidden, errorlist.ErrHandlerNotAllowed)
		}

		return next(c)
	}
}
