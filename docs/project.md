# ���️ 프로젝트 구조 (Project Structure)

이 프로젝트는 **Go (Golang)** 언어와 **Fiber** 웹 프레임워크를 기반으로 하며, 확장성과 유지보수성을 고려한 **Clean Architecture** 및 **DDD(Domain-Driven Design)** 원칙을 따릅니다.

## ��� 디렉토리 구조 (Directory Layout)

```
.
├── cmd/
│   └── myapp/          # 애플리케이션 진입점 (main.go)
├── configs/            # 환경 설정 파일 (local, dev, prod 등)
├── docs/               # 프로젝트 문서 (project, develop 등)
├── infra/              # 인프라 관련 설정 (Docker Compose 등)
├── internal/           # 비공개 애플리케이션 코드 (외부에서 import 불가)
│   ├── config/         # 설정 로드 및 관리 (cleanenv)
│   ├── database/       # 데이터베이스 연결 및 관리
│   ├── feature/        # 도메인별 기능 구현 (Auth, User 등)
│   ├── metrics/        # Prometheus 메트릭 등록 및 초기화
│   ├── middleware/     # Fiber 미들웨어 (CORS, Logger, Recovery)
│   ├── observability/  # OpenTelemetry (Tracing) 관련 설정
│   ├── router/         # 최상위 라우터 설정
│   └── shared/         # 공통적으로 사용되는 모델 및 유틸리티
├── migrations/         # DB 마이그레이션 파일
├── pkg/                # 외부 프로젝트에서도 사용 가능한 라이브러리 코드
│   ├── dbmetrics/      # Prometheus 메트릭 정의 (전역 사용)
│   └── log/            # Context 기반 로깅 (Zap + OTEL Trace ID 주입)
├── scripts/            # 유틸리티 스크립트
└── test/               # 통합 테스트 및 테스트 유틸리티
```

---

## ���️ 아키텍처 및 주요 컴포넌트

### 1. `internal/feature` (도메인 계층)
각 기능(Feature)은 독립적인 패키지로 구성되어 높은 응집도를 가집니다.
- **Handler**: HTTP 요청 처리, 파라미터 파싱 및 검증.
- **Service**: 핵심 비즈니스 로직 수행.
- **Repository**: 데이터베이스와의 직접적인 상호작용 (Store 패턴).
- **Model (DTO/Entity)**: 데이터 구조 정의.

### 2. `pkg` (공용 패키지)
애플리케이션 비즈니스 로직과 무관하게 재사용 가능한 일반적인 기능을 담습니다.
- **dbmetrics**: 애플리케이션 전역에서 사용되는 Prometheus 메트릭을 정의합니다.
- **log**: `zap` 기반 로깅 래퍼이며, `context.Context`에서 Trace ID를 추출해 자동으로 로그에 포함시키는 `InfoCtx`, `ErrorCtx` 함수를 제공합니다.

### 3. `internal/config` (설정 관리)
- `ilyakaznacheev/cleanenv`를 사용하여 `.env` 파일 및 환경변수를 구조체로 로드합니다.

### 4. `internal/middleware` (미들웨어)
- **ApiMetrics**: 요청 처리 시간을 측정하고 Prometheus 메트릭을 업데이트합니다.
- **Recovery**: 패닉 발생 시 서버 중단 방지.

### 5. `internal/observability` & `metrics` (관측성)
- **OpenTelemetry (OTEL)**: 분산 트레이싱을 통해 요청 흐름 추적.
- **Prometheus**: `pkg/dbmetrics`에 정의된 메트릭을 수집하여 노출 (/metrics).

---

## ���️ 기술 스택 (Tech Stack)

| 구분 | 기술 | 설명 |
|---|---|---|
| **Language** | Go 1.25+ | 고성능 시스템 프로그래밍 언어 |
| **Framework** | Fiber v2 | Express 스타일의 고성능 웹 프레임워크 |
| **Database** | PostgreSQL (PGX v5) | 강력한 오픈소스 RDBMS 및 드라이버 |
| **Logging** | Zap (Custom Ctx) | Context 인식 및 Trace ID 자동 주입 로깅 |
| **Config** | Cleanenv | 깔끔한 환경변수 설정 관리 |
| **Observability** | OpenTelemetry | 표준화된 트레이싱 및 모니터링 |
| **Metrics** | Prometheus | 시스템 메트릭 수집 및 시각화 준비 |
