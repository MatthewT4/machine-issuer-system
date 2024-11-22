package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"machineIssuerSystem/internal/model"
)

func (h *handlers) SignUp(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	request := model.SignUpRequest{}
	if err = json.Unmarshal(body, &request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := h.core.SignUp(ctx.Request().Context(), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ctx.SetCookie(&http.Cookie{
		Name:    h.cfg.Auth.CookieName,
		Value:   token,
		Expires: time.Now().Add(time.Duration(h.cfg.Auth.TTL) * time.Hour),
	})

	return ctx.String(http.StatusCreated, "Sign up successfully")
}

func (h *handlers) SignIn(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	request := model.SignInRequest{}

	if err = json.Unmarshal(body, &request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := h.core.SignIn(ctx.Request().Context(), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ctx.SetCookie(&http.Cookie{
		Name:    h.cfg.Auth.CookieName,
		Value:   token,
		Expires: time.Now().Add(time.Duration(h.cfg.Auth.TTL) * time.Hour),
	})

	return ctx.String(http.StatusCreated, "Sign in successfully")
}

func (h *handlers) SignOut(ctx echo.Context) error {
	ctx.SetCookie(&http.Cookie{
		Name:    h.cfg.Auth.CookieName,
		Value:   "",
		Expires: time.Now(),
	})

	return ctx.String(http.StatusOK, "Sign out successfully")
}
