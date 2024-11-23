package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"machineIssuerSystem/internal/model"
	"machineIssuerSystem/pkg/errorlist"
	"machineIssuerSystem/pkg/jwt"
	"machineIssuerSystem/pkg/roles"
)

func (c *Core) SignUp(ctx context.Context, params model.SignUpRequest) (token string, err error) {
	const op = "authCore.SignUp"

	log := c.logger.With(
		slog.String("op", op),
		slog.String("email", params.Email),
	)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err.Error())

		return token, err
	}

	user, err := c.storage.CreateUser(ctx, model.User{
		UUID:         uuid.New(),
		Username:     params.Username,
		Email:        params.Email,
		HashPassword: string(passHash),
		Role:         roles.User,
	})
	if err != nil {
		log.Error("failed to create user", err.Error())

		return token, err
	}

	log.Info("successfully created user", user)

	token, err = jwt.NewToken(user, c.cfg.AuthSecretKey, time.Duration(c.cfg.AuthTTL)*time.Hour)
	if err != nil {
		log.Error("failed to generate token", err.Error())

		return token, err
	}

	return token, nil
}

func (c *Core) SignIn(ctx context.Context, params model.SignInRequest) (token string, err error) {
	const op = "authCore.SignIn"
	ff, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Println(string(ff))
	log := c.logger.With(
		slog.String("op", op),
		slog.String("username", params.Username),
	)

	log.Info("login for user")

	user, err := c.storage.GetUserByUsername(ctx, params.Username)
	if err != nil {
		log.Error("failed to fetch user", err.Error())

		return token, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(params.Password)); err != nil {
		log.Info("invalid password", err.Error())

		return "", fmt.Errorf("%s: %w", op, errorlist.ErrInvalidCredentials)
	}

	token, err = jwt.NewToken(user, c.cfg.AuthSecretKey, time.Duration(c.cfg.AuthTTL)*time.Hour)
	if err != nil {
		log.Error("failed to generate token", err.Error())

		return token, err
	}

	return token, nil
}

func (c *Core) GetPermissionHandler(
	ctx context.Context,
	params model.GetPermissionHandlerRequest,
) (response model.PermissionHandler, err error) {
	response, err = c.storage.GetPermissionHandler(ctx, params)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return response, err
	}

	return response, nil
}
