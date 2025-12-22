package auth

import (
	"context"

	"study/internal/feature/member"
	"study/internal/shared/db"
	"study/internal/shared/model"
	"study/pkg/log"
	"study/pkg/util"
)

// AuthService
// - 인증/회원가입 비즈니스 로직 전용
type AuthService struct {
	MemberRepository *member.MemberRepository
	JwtService       *JwtService
	DB               db.DB
}

// 로그인 요청 DTO
type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe" default:"false"`
}

// 멤버 전달 객체
type MemberResponse struct {
	ID      int64    `json:"id"`
	Email   string   `json:"email"`
	Profile *string  `json:"profile"`
	Roles   []string `json:"roles"`
}

// 로그인 응답 DTO
type LoginResponse struct {
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	Member       MemberResponse `json:"member"`
}

// 생성자
func NewAuthService(MemberRepository *member.MemberRepository, JwtService *JwtService, DB db.DB) *AuthService {
	return &AuthService{MemberRepository: MemberRepository, JwtService: JwtService, DB: DB}
}

// 회원가입
func (s *AuthService) Register(ctx context.Context, member *member.Member) error {

	// 트랜젝션 시작
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 이메일 중복 체크
	_, err = s.MemberRepository.FindByEmail(ctx, tx, member.Email)
	if err == nil {
		return ErrEmailAlreadyExists
	}

	// 비밀번호 암호화
	hashed, err := util.HashString(member.Password)
	if err != nil {
		return err
	}
	member.Password = hashed

	// 상태값 부여
	member.Status = model.StatusActive

	// 회원생성
	err = s.MemberRepository.Create(ctx, tx, member)
	if err != nil {
		return err
	}

	// 기본권한 추가
	err = s.MemberRepository.InsertDefaultRole(ctx, tx, *member.ID)
	if err != nil {
		return err
	}

	// 커밋
	return tx.Commit(ctx)
}

// 로그인
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 이메일로 회원 조회
	member, err := s.MemberRepository.FindByEmail(ctx, s.DB, req.Email)
	if err != nil {
		return nil, ErrInvalidCredential
	}
	log.Info("이메일 조회 성공")

	// 비밀번호 비교
	err = util.VerifyHashString(req.Password, member.Password)
	if err != nil {
		return nil, ErrInvalidCredential
	}
	log.Info("비밀번호 비교 성공")

	// 권한 조회
	roles, err := s.MemberRepository.GetRolesByMemberID(ctx, s.DB, *member.ID)
	if err != nil {
		return nil, err
	}

	// 토큰 생성
	loginResponse, err := s.JwtService.Login(*member.ID)
	if err != nil {
		return nil, err
	}

	loginResponse.Member = MemberResponse{
		ID:      *member.ID,
		Email:   member.Email,
		Profile: member.Profile,
		Roles:   roles,
	}

	return loginResponse, nil
}

// 리프레쉬 토큰으로 로그인 상태 유지
func (s *AuthService) refresh(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// refresh token 검증
	claims, err := s.JwtService.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 회원 찾기
	member, err := s.MemberRepository.FindByID(ctx, s.DB, claims.MemberID)
	if err != nil {
		return nil, err
	}

	// 권한 조회
	roles, err := s.MemberRepository.GetRolesByMemberID(ctx, s.DB, *member.ID)
	if err != nil {
		return nil, err
	}

	// 엑세스 토큰 생성
	accessToken, err := s.JwtService.GenerateAccessToken(*member.ID)
	if err != nil {
		return nil, err
	}

	// 리프레쉬 토큰 제외한 로그인 응답 반환
	return &LoginResponse{
		AccessToken: accessToken,
		Member: MemberResponse{
			ID:      *member.ID,
			Email:   member.Email,
			Profile: member.Profile,
			Roles:   roles,
		},
	}, nil
}
