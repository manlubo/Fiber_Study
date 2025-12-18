# 프로젝트 구조 및 아키텍처 (Project Structure & Architecture)

## 1. 개요

이 프로젝트는 **Go (Fiber)** 기반의 고성능 웹 애플리케이션으로, **도메인 주도 설계(Domain-Driven Design)** 철학을 차용하여 기능별로 높은 응집도를 갖도록 구성되었습니다. `Antigravity` 스타일에 맞춰 단순함, 명확함, 그리고 일관성을 핵심 가치로 삼습니다.

## 2. 기술 스택 (Tech Stack)

- **Language**: Go 1.25+
- **Framework**: [Fiber v2](https://github.com/gofiber/fiber) (Express 스타일의 고성능 웹 프레임워크)
- **Database**: PostgreSQL (Driver: [pgx/v5](https://github.com/jackc/pgx) - High performance connection pool)
- **Configuration**: [cleanenv](https://github.com/ilyakaznacheev/cleanenv) (Struct-based config management)
- **Logging**: [zap](https://github.com/uber-go/zap) (Structured, High-performance logging)

---

## 3. 디렉토리 구조 (Directory Structure)

```
.
├── cmd/
│   └── myapp/              # 애플리케이션 엔트리포인트 (main.go)
├── configs/                # 환경별 설정 파일 (application-dev.yml, prod.yml)
├── docs/                   # 프로젝트 문서 (project.md, develop.md 등)
├── internal/               # 외부에서 import 불가능한 비공개 패키지
│   ├── config/             # 설정 로드 (cleanenv, 구조체 정의)
│   ├── database/           # DB 연결 관리 (pgxpool)
│   ├── feature/            # ⭐ 핵심: 도메인(기능) 단위 패키지 구성
│   │   ├── auth/           # 인증 도메인 (Handler, Service, Router, DTO)
│   │   └── member/         # 회원 도메인 (Repository, Entity)
│   ├── middleware/         # Fiber 글로벌 미들웨어 (CORS, Metrics)
│   ├── router/             # 루트 라우터 및 라우트 등록
│   └── shared/             # 도메인 간 공유되는 모델/코드 (BaseModel 등)
├── pkg/                    # 외부에서도 사용 가능한 범용 패키지
│   ├── dbmetrics/          # DB connection pool 모니터링
│   ├── log/                # Zap 기반 커스텀 로거 (JSON/Console 지원)
│   ├── response/           # API 표준 응답 포맷 (JSend 스타일)
│   └── util/               # 유틸리티 (Path, Hash 등)
├── scripts/                # 보조 스크립트
├── test/                   # 테스트 코드
├── .air.toml               # Air (Live Reload) 설정
├── .env                    # 로컬 개발용 환경변수
└── go.mod                  # Go 모듈 의존성 정의
```

---

## 4. 아키텍처 및 계층 (Layered Architecture)

본 프로젝트는 **3-Layer Architecture**를 기반으로 하지만, 물리적 구조는 **Feature-based Packaging**을 따릅니다.

### 🔄 데이터 흐름 (Data Flow)

`Request` ➡️ **Middleware** ➡️ **Handler** ➡️ **Service** ➡️ **Repository** ➡️ **Database**

### 1️⃣ Handler (Presentation Layer)

- **위치**: `internal/feature/*/handler.go`
- **역할**:
  - HTTP 요청 수신, 파라미터 파싱 및 검증(Validation)
  - `pkg/response`를 사용한 표준 응답 반환
  - **비즈니스 로직을 포함하지 않음** (Service 호출만 수행)

### 2️⃣ Service (Business Layer)

- **위치**: `internal/feature/*/service.go`
- **역할**:
  - 핵심 비즈니스 로직 수행
  - 트랜잭션 단위 관리
  - 여러 Repository를 조합하여 기능 구현

### 3️⃣ Repository (Data Access Layer)

- **위치**: `internal/feature/*/repo.go`
- **역할**:
  - 데이터베이스 직접 접근 (`pgxpool` 사용)
  - 순수 CRUD 쿼리 실행
  - 도메인 로직 포함 금지

---

## 5. 주요 모듈 상세 (Key Components)

### ⚙️ 설정 관리 (`internal/config`)

- `cleanenv` 라이브러리를 사용해 **YAML 파일** + **환경변수(.env)**를 `Config` 구조체 하나로 매핑합니다.
- 실행 환경(`APP_ENV`)에 따라 적절한 `application-{env}.yml`을 로드합니다.
- `Config` 구조체는 `config_struct.go`에 정의되어 있으며, 태그(`yaml`, `env`)를 통해 매핑 규칙을 명시합니다.

### 📝 로깅 시스템 (`pkg/log`)

- **Uber Zap** 기반의 고성능 로거를 래핑(Wrapping)하여 사용합니다.
- **환경별 동작**:
  - `dev`: 사람이 읽기 쉬운 **Color Console** 포맷
  - `prod`: 기계 수집에 용이한 **JSON** 포맷 (ISO8601 Timestamp)
- **Helper**: `log.MapStr`, `log.MapErr` 등을 제공하여 구조화된 필드 로깅을 쉽게 할 수 있도록 지원합니다.

### 🌐 미들웨어 (`internal/middleware`)

- **CORS**: `internal/middleware/cors.go`에서 설정(`config.Cors`)을 기반으로 허용 도메인/메서드를 제어합니다.
- **Metrics**: API 응답 시간, 슬로우 쿼리 등을 측정하여 로그로 남깁니다.

---

## 6. 개발 컨벤션 (Conventions)

1. **명시적 의존성 주입 (DI)**: `main.go`에서 Config, DB, Service, Handler를 생성하고 연결합니다. 전역 변수 사용을 지양합니다.
2. **에러 처리**: 에러는 `handler` 계층까지 전파(`return err`)하여 중앙에서 처리하거나 로그를 남깁니다. 에러 래핑(`fmt.Errorf("%w", err)`)을 권장합니다.
3. **표준 응답**: 모든 API는 성공 시 `response.OK`, 실패 시 에러를 반환하여 일관된 JSON 구조(`result`, `data`, `message`)를 유지합니다.
