# 개발 가이드 (Development Guide)

이 문서는 프로젝트 개발 시 준수해야 할 규칙과 가이드를 제공합니다.

---

## 🛠 개발 환경 (Environment)

- **Language**: Go 1.22+
- **Framework**: Fiber
- **Database**: PostgreSQL (Gorm or raw SQL)
- **Config**: cleanenv (.env)

---

## 📏 코딩 컨벤션 (Conventions)

### 1. 네이밍 & 구조

- **변수명**: 짧고 직관적으로 (`user`, `req`, `resp`).
- **함수명**: 동사+명사 (`CreateUser`, `FindByID`).
- **파일명**: snake_case 권장 (`user_service.go`).

### 2. 에러 처리

- `error`는 항상 마지막 리턴 값으로 반환합니다.
- 커스텀 에러(`errorx` 패키지)를 사용하여 HTTP 상태 코드와 메시지를 명확히 합니다.
- `panic`은 서버 시작 시 등 치명적인 경우에만 사용합니다.

### 3. 로그 (Logging)

- 구조화된 로그를 사용합니다 (JSON format 등).
- 민감 정보(비밀번호 등)는 절대 로그에 남기지 않습니다.

---

## 📊 모니터링 & 메트릭 (Observability)

### Prometheus Metrics

`internal/metrics` 패키지를 통해 서버의 상태와 성능 지표를 수집합니다.

- **Endpoint**: `/metrics` (Prometheus가 scraping)
- **주요 지표**:
  - HTTP 요청 수, 응답 시간, 에러율
  - 시스템 리소스 (Go 런타임 메트릭)
- **사용법**:
  - `api_metrics` 미들웨어가 자동으로 HTTP 요청 정보를 수집합니다.
  - 커스텀 메트릭이 필요한 경우 `metrics.Register...` 함수를 활용합니다.

---

## 🚀 개선 제안 및 로드맵 (Improvements & Roadmap)

### 1. 테스트 강화 (Testing)

- 현재 `test/` 폴더 내의 테스트 커버리지를 높일 필요가 있습니다.
- **Unit Test**: 비즈니스 로직(`Service`) 위주로 Mocking을 활용하여 작성.
- **Integration Test**: 실제 DB를 연결하거나 Docker Compose를 활용하여 API 엔드포인트 테스트.

### 2. CI/CD 파이프라인

- GitHub Actions 등을 도입하여 PR 생성 시 자동 빌드 및 테스트 수행.
- Linting(`golangci-lint`) 자동화.

### 3. API 문서화 (Docs)

- Swagger(OpenAPI) 등을 도입하여 API 명세를 자동화하거나 체계적으로 관리.
- 주석(`// @Summary`) 기반의 `swag` 도구 활용 고려.

### 4. 로깅 고도화

- Request ID(Trace ID)를 도입하여 요청의 흐름을 추적할 수 있도록 개선 (Fiber Middleware 활용).
