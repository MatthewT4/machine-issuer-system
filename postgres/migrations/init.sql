CREATE TABLE servers
(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    cpu integer NOT NULL,
    memory integer NOT null,
    disk integer not null,
    rent_by UUID
);

CREATE TABLE users
(
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    role INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE permission_handlers
(
    id SERIAL PRIMARY KEY,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    roles INTEGER[]
)
