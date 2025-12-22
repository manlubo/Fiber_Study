# ��� 개발 가이드 (Development Guide)

프로젝트 개발 및 유지보수를 위한 가이드라인과 로드맵입니다.

## ✅ 현재 구현 상태 (Current Status)

- [x] **기본 구조**: Clean Architecture 기반의 폴더 구조 확립.
- [x] **웹 서버**: Fiber v2 프레임워크 적용 및 라우팅 설정.
- [x] **데이터베이스**: PostgreSQL 연결 및 PGX 드라이버 적용.
- [x] **설정 관리**: `cleanenv`를 이용한 환경 변수 관리.
- [x] **관측성 (Observability)**: OpenTelemetry 트레이싱 및 Prometheus 메트릭 엔드포인트 구현.
- [x] **미들웨어**: CORS, 로깅(Zap), 패닉 복구 미들웨어 적용 완료.

---

## ��️ 개선 필요 사항 (Improvements & TODO)

현재 프로젝트의 완성도를 높이기 위해 다음과 같은 작업들이 필요합니다.

### 1. 테스트 강화
- **Unit Test**: `feature` 계층의 Service 로직에 대한 단위 테스트 작성 (Mocking 활용).
- **Integration Test**: 실제 DB와 연동되는 통합 테스트 추가 (`testcontainers-go` 활용 권장).

### 2. API 문서화 (Swagger)
- **Swaggo 도입**: `swaggo/swag`를 사용하여 코드 주석 기반으로 Swagger 문서를 자동 생성.
- **엔드포인트 명세**: 요청/응답 스키마를 명확히 정의하여 프론트엔드 개발자와의 협업 효율 증대.

### 3. CI/CD 파이프라인 (권장)
- **GitHub Actions**: 코드 푸시 시 자동으로 빌드 및 테스트를 수행하는 워크플로우 추가.
- **Linting**: `golangci-lint`를 적용하여 코드 품질 자동 점검.

### 4. 에러 핸들링 고도화
- **전역 에러 처리**: `GlobalErrorHandler`를 사용하여 예상치 못한 에러도 일관된 JSON 포맷으로 응답.
- **Custom Error**: 도메인별 비즈니스 에러 타입 정의 (`ErrUserNotFound`, `ErrInvalidPassword` 등).

### 5. 보안 강화
- **JWT 관리**: 토큰 만료 시간, 리프레시 토큰 로직 재점검 및 보안 강화.
- **Security Headers**: Helmet 등 보안 헤더 미들웨어 추가.

---

## ��� 기여 방법 (Contribution)

1. **이슈 확인**: `docs/develop.md`의 개선 사항 중 작업할 내용을 선택합니다.
2. **브랜치 생성**: `feature/기능명` 또는 `fix/버그명` 형태로 브랜치를 생성합니다.
3. **코드 작성**: 기존 코드 스타일과 컨벤션을 준수하여 코드를 작성합니다.
4. **테스트**: 변경 사항에 대한 테스트를 수행합니다.
5. **PR 요청**: 작업 내용을 상세히 작성하여 Pull Request를 보냅니다.
