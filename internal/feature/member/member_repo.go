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

// 레포지토리 스트럭처
type MemberRepository struct {
	db db.DB
}

// 레포지토리 생성
func NewMemberRepository(db db.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

// 회원가입
func (r *MemberRepository) Create(ctx context.Context, member *Member) error {

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
	err := r.db.QueryRow(
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

// 이메일로 회원 조회
func (r *MemberRepository) FindByEmail(ctx context.Context, email string) (*Member, error) {
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
	err := r.db.QueryRow(ctx, query, email).Scan(
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
func (r *MemberRepository) FindByID(ctx context.Context, id int64) (*Member, error) {
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
	err := r.db.QueryRow(ctx, query, id).Scan(
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
