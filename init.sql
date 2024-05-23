-- Table: public.servers

-- DROP TABLE IF EXISTS public.servers;

CREATE TABLE IF NOT EXISTS public.servers
(
    id bigint NOT NULL DEFAULT nextval('servers_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    ipv4 character varying(15) COLLATE pg_catalog."default" NOT NULL,
    status character varying(50) COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT servers_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.servers
    OWNER to postgres;
-- Index: idx_servers_ipv4

-- DROP INDEX IF EXISTS public.idx_servers_ipv4;

CREATE UNIQUE INDEX IF NOT EXISTS idx_servers_ipv4
    ON public.servers USING btree
    (ipv4 COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: idx_servers_name

-- DROP INDEX IF EXISTS public.idx_servers_name;

CREATE UNIQUE INDEX IF NOT EXISTS idx_servers_name
    ON public.servers USING btree
    (name COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- Table: public.users

-- DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default",
    role text COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;
-- Index: idx_users_username

-- DROP INDEX IF EXISTS public.idx_users_username;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username
    ON public.users USING btree
    (username COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

INSERT INTO public.users (username, password, role) VALUES
('duc', '$2a$10$g7j7FwaE8HeYnGmSD2P9yuLX1sPHx4qG5UofS4F2QTEhivcncZM1m', 'admin')