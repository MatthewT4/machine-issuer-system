package storage

const (
	queryCreateUser = `
INSERT INTO user (username, email, password, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING id, username, email, password, created_at, updated_at
`

	queryGetUserByUsername = `
SELECT id, username, email, password, created_at, updated_at 
FROM user
WHERE username = $1
`
)
