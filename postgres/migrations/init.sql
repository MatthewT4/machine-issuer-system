CREATE TABLE servers
(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    cpu integer NOT NULL,
    memory integer NOT null,
    disk integer not null,
    rent_by UUID,
    ip TEXT NOT NULL,
    rent_until TIMESTAMP
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
       ('POST', '/rent/%', '{1, 2, 3}'),
       ('DELETE', '/rent/%', '{1, 2, 3}'),
       ('GET', '/metrics/%', '{1, 2, 3}'),
       ('GET', '/metrics', '{0, 1, 2, 3}');

INSERT INTO public.servers (id, title, cpu, memory, disk, rent_by, ip)
VALUES ('4f45b109-0e26-43b0-abcb-52b216b69a1e', 'yandex-cloud-1', 12, 8192, 16000, null, '51.250.41.219');

INSERT INTO public.servers (id, title, cpu, memory, disk, rent_by, ip)
VALUES ('4f45b109-0e26-43b0-abcb-52b216b69a1a', 'baikal-1', 12, 8192, 16000, null, '87.251.74.153');


INSERT INTO public.users (id, email, username, password, role, created_at, updated_at) VALUES ('00000000-0000-0000-0000-000000000001', '2342fdsdf', 'test', '$2a$10$sKMAr4MGSD5x.yOX8/wD3OwzuuxkxyxbGTJ7oMrFH1RC.l6J084..', 1, '2024-11-23 21:48:02.000000', '2024-11-23 21:48:05.000000');
