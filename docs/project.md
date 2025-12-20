# 프로젝트 구조 (Project Structure)

이 문서는 프로젝트의 디렉토리 구조와 각 패키지의 역할을 설명합니다.  
본 프로젝트는 **Standard Go Project Layout**을 기반으로 하며, **Domain-Driven Design(DDD)**의 개념을 일부 차용하여 `feature` 단위로 비즈니스 로직을 구성했습니다.

---

## 📂 디렉토리 상세 (Directory Details)

### `cmd/`

애플리케이션의 진입점(Entry Point)입니다.

- **myapp/**: 메인 애플리케이션 실행 파일이 위치합니다. (`main.go`)

### `configs/`

환경 변수 및 설정 파일들이 위치합니다.

- `.env`: 로컬 개발용 환경 변수 (git 제외)
- `config.yaml`: 기본 설정 파일 (필요 시)

### `docs/`

프로젝트 관련 문서들이 위치합니다.

- `project.md`: 프로젝트 구조 설명 (본 문서)
- `develop.md`: 개발 가이드 및 규칙

### `internal/`

외부에서 임포트할 수 없는 비공개 애플리케이션 코드입니다.

#### `internal/config`

- `cleanenv` 등을 사용하여 환경 변수나 설정 파일을 로드하고 구조체로 매핑합니다.

#### `internal/database`

- 데이터베이스 연결 및 초기화 로직을 담당합니다.

#### `internal/feature` (핵심 도메인)

기능(Feature) 단위로 패키지를 구분합니다. 각 기능 폴더 내부는 **Layered Architecture**를 따릅니다.

- **auth/**: 인증/인가 관련 로직 (로그인, 토큰 발급 등)
- **member/**: 회원 관리 로직
- _각 feature 구조_:
  - `handler/`: HTTP 요청 처리 (Controller)
  - `service/`: 비즈니스 로직
  - `repository/`: DB 데이터 접근
  - `dto/`: 데이터 전송 객체
  - `entity/`: 도메인 모델

#### `internal/metrics` (New)

- **Prometheus** 메트릭 수집 및 노출을 담당합니다.
- `exporter.go`: 메트릭 설정 및 등록
- `http.go`: 메트릭 수집을 위한 HTTP 핸들러 및 유틸리티
- `middleware.go`: Gin/Fiber 등 웹 프레임워크와 연동되는 메트릭 미들웨어

#### `internal/middleware`

- **Fiber** 미들웨어 모음입니다.
- `cors.go`: CORS 설정
- `auth.go`: JWT 인증 미들웨어
- `api_metrics.go`: API 요청 메트릭 수집 미들웨어

#### `internal/router`

- 라우팅 등록 및 설정을 담당합니다. 각 `feature`의 핸들러를 엔드포인트에 매핑합니다.

#### `internal/shared`

- 전역에서 공통으로 사용되는 유틸리티 및 모델입니다.
- `errorx/`: 커스텀 에러 처리
- `db/`: DB 관련 공통 유틸 (Pagination 등)
- `model/`: 공통 데이터 모델

### `migrations/`

- 데이터베이스 스키마 마이그레이션 파일들이 위치합니다.

### `pkg/`

- 외부 프로젝트에서도 사용할 수 있는 범용 라이브러리 코드입니다. (현재 프로젝트에서는 내부 로직 위주이므로 사용 빈도가 낮을 수 있음)

### `scripts/`

- 빌드, 배포, 테스팅 등을 위한 유틸리티 스크립트입니다.

### `test/`

- 테스트 코드들이 위치합니다.
- `unit/`: 단위 테스트
- `integration/`: 통합 테스트

---

## 🏗 아키텍처 원칙 (Architecture Principles)

1. **관심사의 분리 (Separation of Concerns)**:

   - `Handler`는 HTTP 요청/응답만 처리합니다.
   - `Service`는 비즈니스 로직만 수행합니다.
   - `Repository`는 DB 쿼리만 수행합니다.

2. **의존성 주입 (Dependency Injection)**:

   - 각 계층은 인터페이스를 통해 협력하며, 초기화 시점에 의존성을 주입받습니다.

3. **기능 중심 패키징 (Package by Feature)**:
   - 관련된 코드를 기능(`internal/feature/xxx`)별로 모아서 응집도를 높입니다.
