BEGIN;

CREATE TABLE IF NOT EXISTS goods
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT             NOT NULL,
    description TEXT             NOT NULL,
    unit        TEXT             NOT NULL,
    cost        DOUBLE PRECISION NOT NULL,
    created_by  BIGINT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    enabled     BOOLEAN                  DEFAULT TRUE,

    CONSTRAINT fk_user FOREIGN KEY (created_by) REFERENCES users (id)
);

END;