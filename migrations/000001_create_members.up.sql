CREATE TABLE members (
    member_id BIGSERIAL PRIMARY KEY,

    email TEXT NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,

    tel TEXT,
    address TEXT,
    profile TEXT,

    status TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX uq_members_email_active_ci
ON members (lower(email))
WHERE deleted_at IS NULL;

