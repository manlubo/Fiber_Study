# ��� 개발 가이드 (Development Guide)

프로젝트 개발 및 유지보수를 위한 가이드라인과 로드맵입니다.

## ✅ 현재 구현 상태 (Current Status)

- [x] **기본 구조**: Clean Architecture 기반의 폴더 구조 확립.
- [x] **웹 서버**: Fiber v2 + PostgreSQL + PGX.
- [x] **관측성**: OpenTelemetry Tracing + Prometheus Metrics.
- [x] **로깅**: Zap 기반의 Context Aware Logging 구현 (Trace ID 자동 주입).

---

## ��� 개발 규칙 & 컨벤션 (Conventions)

### 1. 로깅 (Logging) & 트레이싱 (Tracing)
분산 트레이싱 환경에서 로그와 트레이스를 연결하기 위해 **반드시 Context 기반 로깅**을 사용합니다.

- **규칙**: `context.Context`가 사용 가능한 모든 계층(Handler, Service, Repository)에서는 `log.Info`, `log.Error` 대신 **`log.InfoCtx`, `log.ErrorCtx`** 함수를 사용해야 합니다.
- **이유**: 해당 함수들은 Context 내의 SpanContext에서 `trace_id`와 `span_id`를 자동으로 추출하여 로그 필드에 추가합니다. 이는 Grafana/Tempo 등에서 로그와 트레이스를 연관 짓는 핵심 키가 됩니다.

**예시 (Good)**:
```go
func (s *AuthService) Login(ctx context.Context, req LoginRequest) {
    // ...
    // 자동으로 trace_id가 로그에 남음
    log.InfoCtx(ctx, "로그인 성공", log.String("email", req.Email))
}
```

**예시 (Bad)**:
```go
func (s *AuthService) Login(ctx context.Context, req LoginRequest) {
    // 트레이스 연결 끊김
    log.Info("로그인 성공", log.String("email", req.Email)) 
}
```

### 2. 메트릭 (Metrics)
새로운 메트릭을 추가할 때는 다음 절차를 따릅니다.

1. **정의**: `pkg/dbmetrics` 패키지에 메트릭 변수(,  등)를 정의합니다.
   - 이유: 패키지 순환 참조 방지 및 메트릭 정의의 중앙 관리.
2. **등록**: `internal/metrics/init.go` 내의 `Init()` 함수에서 `prometheus.MustRegister()`를 통해 등록합니다.
3. **사용**: 필요한 미들웨어난 로직에서 `dbmetrics.MyMetric.Inc()` 와 같이 사용합니다.

---

## ���️ 개선 필요 사항 (Improvements & TODO)

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

### 5. 보안 강화
- **JWT 관리**: 토큰 만료 시간, 리프레시 토큰 로직 재점검 및 보안 강화.
- **Security Headers**: Helmet 등 보안 헤더 미들웨어 추가.
