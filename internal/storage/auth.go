package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"

	"machineIssuerSystem/internal/model"
)

func (p *PgStorage) CreateUser(ctx context.Context, user model.User) (result model.User, err error) {
	const op = "storage.CreateUser"

	err = p.connections.QueryRow(
		ctx,
		queryCreateUser,
		user.UUID,
		user.Email,
		user.Username,
		user.HashPassword,
		user.Role,
	).Scan(&result.UUID, &result.Username, &result.Email, &result.HashPassword, &result.Role, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (p *PgStorage) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	const op = "storage.GetUserByUsername"

	rows, err := p.connections.Query(
		ctx,
		queryGetUserByUsername,
		username,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
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
	).Scan(&result.ID, &result.Method, &result.Path, &result.Roles)
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
