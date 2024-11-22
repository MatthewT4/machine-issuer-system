CREATE TABLE users
(
    id UUID PRIMARY KEY,
    login TEXT NOT NULL,
    password TEXT NOT NULL,
    is_admin boolean NOT NULL
);

CREATE TABLE servers
(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    cpu integer NOT NULL,
    memory integer NOT null,
    disk integer not null,
    rent_by UUID
);

