#  개발 가이드 (Development Guide)

프로젝트 개발 및 유지보수를 위한 가이드라인과 로드맵입니다.

## ✅ 현재 구현 상태 (Current Status)

- [x] **기본 구조**: Clean Architecture 기반의 폴더 구조 확립.
- [x] **웹 서버**: Fiber v2 + PostgreSQL + PGX.
- [x] **관측성 (Observability)**: 
  - OpenTelemetry Tracing (Jaeger/Tempo) 적용.
  - Span 자동 생성 및 Context Propagation (Handler -> Service -> Repo).
  - Prometheus Metric Exporter 구현.
- [x] **로깅 (Logging)**: Zap 기반의 Context Aware Logging (Trace ID 자동 주입).

---

##  개발 규칙 & 컨벤션 요약 (Conventions Summary)

> 자세한 내용은 [convention.md](./convention.md)를 참고하세요.

### 1. 로깅 (Logging)
- **반드시 Context 포함**: `log.InfoCtx(ctx, ...)` 사용.
- 개인정보(PII) 기록 금지.

### 2. 에러 처리 및 트레이싱
- **Service Layer**:
  - `RecordServiceError`: 시스템 장애로 이어지는 에러 (Trace Error).
  - `RecordBusinessError`: 단순 비즈니스 로직 실패 (Trace Success, Event Log).
- **Attributes**: 검색을 위한 ID, Type, Status는 적극 권장하되, PII나 보안 토큰은 절대 남기지 않습니다.

---

## ️ 개선 필요 사항 (Improvements & TODO)

현재 프로젝트의 완성도를 높이기 위해 다음과 같은 작업들이 필요합니다.

### 1. 테스트 강화
- **Unit Test**: `feature` 계층의 Service 로직에 대한 단위 테스트 작성.
- **Integration Test**: `testcontainers-go`를 활용한 DB 통합 테스트.

### 2. API 문서화
- **Swagger**: `swaggo` 도입하여 API 명세 자동화.

### 3. CI/CD 파이프라인
- GitHub Actions를 통한 빌드/테스트/린트 자동화.

### 4. 보안 강화
- JWT 만료/갱신 로직 고도화 및 보안 헤더(Helmet) 적용.
