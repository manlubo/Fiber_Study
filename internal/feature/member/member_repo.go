package member

import (
	"context"

	"study/internal/shared/db"
	"study/internal/shared/model"
)

// DB 테이블과 일치
type Member struct {
	model.BaseModel
	ID       *int64  `db:"member_id" json:"memberId"`
	Email    string  `db:"email" json:"email"`
	Password string  `db:"password" json:"password"`
	Name     string  `db:"name" json:"name"`
	Tel      *string `db:"tel" json:"tel"`
	Address  *string `db:"address" json:"address"`
	Profile  *string `db:"profile" json:"profile"`
}

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

// 레포지토리 스트럭처
type MemberRepository struct{}

// 레포지토리 생성
func NewMemberRepository() *MemberRepository {
	return &MemberRepository{}
}

// 회원가입
func (r *MemberRepository) Create(ctx context.Context, exec db.Execer, member *Member) error {
	query := `
		INSERT INTO members (
			email,
			password,
			name,
			tel,
			address,
			profile,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING member_id
	`

	var id int64
	err := exec.QueryRow(
		ctx,
		query,
		member.Email,
		member.Password,
		member.Name,
		member.Tel,
		member.Address,
		member.Profile,
		member.Status,
	).Scan(&id)

	if err != nil {
		return err
	}

	member.ID = &id
	return nil
}

// 권한 추가
func (r *MemberRepository) InsertRole(ctx context.Context, exec db.Execer, memberID int64, role Role) error {
	query := `
		INSERT INTO member_roles (member_id, role)
		VALUES ($1, $2)
	`

	_, err := exec.Exec(ctx, query, memberID, role)
	return err
}

// 권한 조회
func (r *MemberRepository) GetRolesByMemberID(
	ctx context.Context,
	exec db.Execer,
	memberID int64,
) ([]Role, error) {
	query := `
		SELECT role
		FROM member_roles
		WHERE member_id = $1
	`

	rows, err := exec.Query(ctx, query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]Role, 0)
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// 이메일로 회원 조회
func (r *MemberRepository) FindByEmail(ctx context.Context, exec db.Execer, email string) (*Member, error) {
	query := `
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
		WHERE email = $1
	`

	var member Member
	err := exec.QueryRow(ctx, query, email).Scan(
		&member.ID,
		&member.Email,
		&member.Password,
		&member.Name,
		&member.Tel,
		&member.Address,
		&member.Profile,
		&member.Status,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.DeletedAt,
	)
	// 멤버가 없을 경우
	if err != nil {
		return nil, err
	}

	return &member, nil
}

// ID로 회원 조회
func (r *MemberRepository) FindByID(ctx context.Context, exec db.Execer, id int64) (*Member, error) {
	query := `
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
		WHERE member_id = $1
	`

	var member Member
	err := exec.QueryRow(ctx, query, id).Scan(
		&member.ID,
		&member.Email,
		&member.Password,
		&member.Name,
		&member.Tel,
		&member.Address,
		&member.Profile,
		&member.Status,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.DeletedAt,
	)
	// 멤버가 없을 경우
	if err != nil {
		return nil, err
	}

	return &member, nil
}
