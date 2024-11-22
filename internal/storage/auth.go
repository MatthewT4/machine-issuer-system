package storage

import (
	"context"
	"fmt"

	"machineIssuerSystem/internal/model"
)

func (p *PgStorage) CreateUser(ctx context.Context, user model.User) (result model.User, err error) {
	const op = "storage.CreateUser"

	err = p.connections.QueryRow(
		ctx,
		queryCreateUser,
		user.Username,
		user.Email,
		user.HashPassword,
	).Scan(&result)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (p *PgStorage) GetUserByUsername(ctx context.Context, username string) (result model.User, err error) {
	const op = "storage.GetUserByUsername"

	err = p.connections.QueryRow(
		ctx,
		queryGetUserByUsername,
		username,
	).Scan(&result)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (p *PgStorage) GetPermissionHandler(
	ctx context.Context,
	params model.GetPermissionHandlerRequest,
) (result model.PermissionHandler, err error) {
	const op = "storage.GetPermissionHandler"

	err = p.connections.QueryRow(
		ctx,
		queryGetPermissionHandler,
		params.Method,
		params.Path,
	).Scan(&result)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
