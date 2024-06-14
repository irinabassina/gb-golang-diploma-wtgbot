BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id         BIGINT  NOT NULL,
    name       TEXT    NOT NULL,
    role       TEXT    NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    enabled    BOOLEAN NOT NULL         DEFAULT TRUE,

    CONSTRAINT pk_users_idx PRIMARY KEY (id)
);

END;