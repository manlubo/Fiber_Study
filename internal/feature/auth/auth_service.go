package auth

import (
	"context"

	"study/internal/feature/member"
	"study/internal/observability"
	"study/internal/shared/db"
	"study/internal/shared/model"
	"study/pkg/log"
	"study/pkg/util"

	"go.opentelemetry.io/otel/attribute"
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
	ID      int64         `json:"id"`
	Email   string        `json:"email"`
	Profile *string       `json:"profile"`
	Roles   []member.Role `json:"roles"`
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
func (s *AuthService) Register(ctx context.Context, m *member.Member) (err error) {
	ctx, span, start := observability.StartServiceSpan(ctx, "Register")
	defer observability.EndSpanWithLatency(span, start, 100)

	// 트랜젝션 시작
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		observability.RecordServiceError(span, err)
		return err
	}
	defer tx.Rollback(ctx)

	// 이메일 중복 체크
	_, err = s.MemberRepository.FindByEmail(ctx, tx, m.Email)
	if err == nil {
		observability.RecordBusinessError(span, ErrEmailAlreadyExists)
		return ErrEmailAlreadyExists
	}

	// 비밀번호 암호화
	hashed, err := util.HashString(m.Password)
	if err != nil {
		observability.RecordServiceError(span, err)
		return err
	}
	m.Password = hashed

	// 상태값 부여
	m.Status = model.StatusActive

	// 회원생성
	err = s.MemberRepository.Create(ctx, tx, m)
	if err != nil {
		observability.RecordServiceError(span, err)
		return err
	}

	// 기본권한 추가
	err = s.MemberRepository.InsertRole(ctx, tx, *m.ID, member.RoleUser)
	if err != nil {
		observability.RecordServiceError(span, err)
		return err
	}

	span.SetAttributes(
		attribute.String("auth.type", "register"),
		attribute.Int64("member.id", *m.ID),
	)

	log.InfoCtx(ctx, "회원가입 성공")

	// 커밋
	if err = tx.Commit(ctx); err != nil {
		observability.RecordServiceError(span, err)
		return err
	}

	return nil
}

// 로그인
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (resp *LoginResponse, err error) {
	ctx, span, start := observability.StartServiceSpan(ctx, "Login")
	defer observability.EndSpanWithLatency(span, start, 0)

	// 이메일로 회원 조회
	member, err := s.MemberRepository.FindByEmail(ctx, s.DB, req.Email)
	if err != nil {
		observability.RecordBusinessError(span, ErrInvalidCredential)
		return nil, ErrInvalidCredential
	}

	// 비밀번호 비교
	err = util.VerifyHashString(req.Password, member.Password)
	if err != nil {
		observability.RecordBusinessError(span, ErrInvalidCredential)
		return nil, ErrInvalidCredential
	}

	// 권한 조회
	roles, err := s.MemberRepository.GetRolesByMemberID(ctx, s.DB, *member.ID)
	if err != nil {
		observability.RecordServiceError(span, err)
		return nil, err
	}

	// 토큰 생성
	loginResponse, err := s.JwtService.Login(*member.ID)
	if err != nil {
		observability.RecordServiceError(span, err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("auth.type", "login"),
		attribute.Int64("member.id", *member.ID),
	)

	loginResponse.Member = MemberResponse{
		ID:      *member.ID,
		Email:   member.Email,
		Profile: member.Profile,
		Roles:   roles,
	}

	log.InfoCtx(ctx, "로그인 성공")
	return loginResponse, nil
}

// 리프레쉬 토큰으로 로그인 상태 유지
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (resp *LoginResponse, err error) {
	ctx, span, start := observability.StartServiceSpan(ctx, "Refresh")
	defer observability.EndSpanWithLatency(span, start, 30)

	// refresh token 검증
	claims, err := s.JwtService.VerifyRefreshToken(refreshToken)
	if err != nil {

		switch err {
		case ErrTokenExpired, ErrTokenTypeWrong:
			// 세션만료, 잘못된 토큰 사용
			observability.RecordBusinessError(span, err)
			return nil, err

		default:
			// 위조, 서명 오류, 내부 파싱 문제
			observability.RecordServiceError(span, err)
			return nil, err
		}
	}

	// 회원 찾기
	member, err := s.MemberRepository.FindByID(ctx, s.DB, claims.MemberID)
	if err != nil {
		observability.RecordServiceError(span, err)
		return nil, err
	}

	// 권한 조회
	roles, err := s.MemberRepository.GetRolesByMemberID(ctx, s.DB, *member.ID)
	if err != nil {
		observability.RecordServiceError(span, err)
		return nil, err
	}

	// 엑세스 토큰 생성
	accessToken, err := s.JwtService.GenerateAccessToken(*member.ID)
	if err != nil {
		observability.RecordServiceError(span, err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("auth.type", "refresh"),
		attribute.Int64("member.id", *member.ID),
	)

	log.InfoCtx(ctx, "로그인 유지 성공")
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
