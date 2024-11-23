CREATE TABLE servers
(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    cpu integer NOT NULL,
    memory integer NOT null,
    disk integer not null,
    rent_by UUID,
    ip TEXT NOT NULL
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
);

INSERT INTO permission_handlers (method, path, roles)
VALUES ('POST', '/auth/signup', '{0, 1, 2, 3}'),
       ('POST', '/auth/signin', '{0, 1, 2, 3}'),
       ('GET', '/servers/available', '{1, 2, 3}'),
       ('GET', '/servers/my', '{1, 2, 3}'),
       ('POST', '/rent/{server_id}', '{1, 2, 3}'),
       ('DELETE', 'rent/{server_id}', '{1, 2, 3}'),
       ('GET', 'metrics/{server_id}', '{1, 2, 3}');