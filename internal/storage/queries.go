package storage

const (
	queryCreateUser = `
INSERT INTO users (id, email, username, password, role, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING id, username, email, password, role, created_at, updated_at
`

	queryGetUserByUsername = `
SELECT id, email, username, password, role, created_at, updated_at 
FROM users
WHERE username = $1
`

	queryGetUserByID = `
SELECT id, email, username, password, role, created_at, updated_at 
FROM users
WHERE id = $1
`

	queryGetPermissionHandler = `
SELECT id, method, path, roles
FROM permission_handlers
WHERE method = $1 and $2 ILIKE path;
`

	queryFetchExpiredServers = `
SELECT id, title, cpu, memory, disk, rent_by, ip, rent_until
FROM servers
WHERE rent_until < NOW() + INTERVAL '3 hours'`
)
