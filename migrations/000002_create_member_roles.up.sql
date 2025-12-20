CREATE TABLE member_roles (
	member_role_id BIGSERIAL PRIMARY KEY,
	member_id BIGINT NOT NULL,
	role VARCHAR(20) NOT NULL,
	UNIQUE (member_id, role),
	CONSTRAINT fk_member_roles_member
		FOREIGN KEY (member_id)
		REFERENCES members(member_id)
		ON DELETE CASCADE
);
