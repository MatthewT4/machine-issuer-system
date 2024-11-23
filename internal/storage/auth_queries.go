package storage

const (
	queryCreateUser = `
INSERT INTO users (id, email, username, password, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING id, username, email, password, created_at, updated_at
`

	queryGetUserByUsername = `
SELECT id, email, username, password, role, created_at, updated_at 
FROM users
WHERE email = $1
`

	queryGetPermissionHandler = `
SELECT id, method, path, roles
FROM permission_handlers
WHERE method = $1 and path = $2
`
)
