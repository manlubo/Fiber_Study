-- name: CreateMember :one
INSERT INTO members (
    email,
    password,
    name,
    status
) VALUES (
    $1, $2, $3, $4
)
RETURNING member_id;


-- name: InsertMemberRole :exec
INSERT INTO member_roles (
    member_id,
    role
) VALUES (
    $1, $2
);


-- name: GetRolesByMemberID :many
SELECT role
FROM member_roles
WHERE member_id = $1;


-- name: FindMemberByEmail :one
SELECT
    member_id,
    email,
    password,
    name,
    tel,
    address,
    profile,
    status,
    created_at,
    updated_at,
    deleted_at
FROM members
WHERE email = $1;


-- name: FindMemberByID :one
SELECT
    member_id,
    email,
    password,
    name,
    tel,
    address,
    profile,
    status,
    created_at,
    updated_at,
    deleted_at
FROM members
WHERE member_id = $1;
