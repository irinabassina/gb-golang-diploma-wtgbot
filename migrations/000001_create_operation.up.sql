BEGIN;

CREATE TABLE IF NOT EXISTS operations
(
    id              BIGSERIAL PRIMARY KEY,
    category_id     bigserial        NOT NULL,
    value           DOUBLE PRECISION NOT NULL,
    current_balance DOUBLE PRECISION NOT NULL,
    created_by      BIGINT,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_user FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES goods (id)
);

END;