package storage

const (
	queryCreateUser = `
INSERT INTO user (id, email, username, password, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING id, username, email, password, created_at, updated_at
`

	queryGetUserByUsername = `
SELECT id, email, username, password, role, created_at, updated_at 
FROM user
WHERE username = $1
`
)
