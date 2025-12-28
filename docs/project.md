# 프로젝트 개요

이 문서는 프로젝트의 전체 구조와 아키텍처, 그리고 주요 기술적 의사결정에 대해 설명합니다.  
본 프로젝트는 **Go (Fiber)** 기반의 웹 애플리케이션으로, **sqlc**를 활용한 타입 안전한 데이터 접근과 **도메인 중심(Feature-based)** 디렉토리 구조를 따르고 있습니다.

## 1. 전체 디렉토리 구조

프로젝트는 Go 표준 레이아웃(Standard Go Layout)을 기반으로 하되, 실용적인 도메인 분리 방식을 채택하고 있습니다.

```
.
├── cmd/
│   └── myapp/             # 애플리케이션 엔트리 포인트 (main.go)
├── internal/
│   ├── config/            # 설정 로드 및 관리 (env, yaml)
│   ├── database/          # DB 연결 설정 (pgx connection pool)
│   ├── feature/           # 도메인 별 비즈니스 로직 (Handler + Service)
│   │   ├── auth/          # 예: 인증 도메인
│   │   └── member/        # 예: 회원 도메인
│   ├── query/             # sqlc로 생성된 DB 접근 코드 (Repository 역할)
│   ├── router/            # 라우터 등록 및 의존성 주입 (Wiring)
│   ├── middleware/        # 공통 미들웨어 (CORS, Auth 등)
│   ├── observability/     # Tracing (OpenTelemetry) 관련 설정
│   └── shared/            # 공용 유틸리티, 모델, 에러 정의
├── docs/                  # 프로젝트 문서
├── sqlc.yaml              # sqlc 설정 파일
└── go.mod                 # 의존성 관리
```

---

## 2. 레이어별 역할 및 책임

이 프로젝트는 **Handler → Service → Query(Data Access)** 의 3계층 구조를 따릅니다.

### 2.1. Presentation Layer (Handler)

- **위치**: `internal/feature/{domain}/*_handler.go`
- **역할**: HTTP 요청/응답 처리
- **책임**:
  - Request Body를 **DTO**(Data Transfer Object)로 파싱 및 유효성 검증
  - Service 계층 호출 (DTO 전달)
  - 적절한 HTTP Status 및 JSON 응답 반환
  - **비즈니스 로직 및 DB Entity 직접 참조 금지**

### 2.2. Business Layer (Service)

- **위치**: `internal/feature/{domain}/*_service.go`
- **역할**: 핵심 비즈니스 로직 수행
- **책임**:
  - 트랜잭션 관리 (`pool.Begin`, `queries.WithTx`)
  - **DTO**를 받아 비즈니스 로직 수행 후, 필요 시 도메인 모델 또는 응답 DTO 반환
  - 도메인 규칙 검증 (예: 이메일 중복 체크, 비밀번호 해싱)
  - Query 계층을 조합하여 기능 수행
  - Observability (Tracing) 스팬 생성 및 에러 기록

### 2.3. Data Access Layer (Query)

- **위치**: `internal/query`
- **역할**: 데이터베이스 상호작용
- **특징**:
  - **sqlc**를 통해 SQL 파일로부터 자동 생성된 Go 코드를 사용
  - `Queries` 구조체를 통해 DB 메서드 제공 (`FindMemberByEmail`, `CreateMember` 등)
  - ORM을 사용하지 않고 Raw SQL에 가까운 명확성을 유지하며 타입 안전성 확보

### 2.4. Infrastructure (Config & Database)

- **Config**: `.env` 파일과 환경 변수를 로드하여 애플리케이션 설정 구조체(`cfg`)로 매핑
- **Database**: `pgx/v5` 드라이버를 사용하여 PostgreSQL 연결 풀(`pgxpool`) 관리

---

## 3. 주요 흐름 (Workflow)

요청이 들어왔을 때의 처리 흐름은 다음과 같습니다.

1.  **Entry (`cmd/myapp/main.go`)**:
    - Config 로드, DB 연결, Logger/Tracer 초기화
    - `query.New(db)`로 Query 인스턴스 생성
    - Service 생성 및 Service에 Query 주입
    - Router에 Handler 등록 (DI 수행)
2.  **Request**:
    - Client → Middleware (Logging, CORS, Auth) → Handler
3.  **Processing**:
    - **Handler**: 요청 파싱 (DTO 바인딩) → **Service** 호출
    - **Service**: 트랜잭션 시작 → **Query** 메서드 호출 (비즈니스 로직 수행) → 결과 반환
4.  **Data Access**:
    - **Query**: Generated Code → `pgx` Driver → PostgreSQL
5.  **Response**:
    - 결과가 역순으로 전달되어 Handler가 최종 JSON 응답

---

## 4. 기술 스택 및 선정 이유

| 기술                 | 선정 이유                                                                                                 |
| -------------------- | --------------------------------------------------------------------------------------------------------- |
| **Go & Fiber**       | 높은 성능과 간결한 동시성 처리, Express.js와 유사한 쉬운 라우팅 API                                       |
| **PostgreSQL & pgx** | 신뢰성 높은 RDBMS와 Go 생태계에서 가장 성능이 좋은 드라이버 사용                                          |
| **sqlc**             | SQL을 직접 작성하여 쿼리 최적화가 용이하면서도, 컴파일 타임에 타입을 보장받을 수 있음 (ORM의 복잡성 제거) |
| **OpenTelemetry**    | 분산 환경에서의 트레이싱 및 모니터링 표준 준수                                                            |

## 5. 설계 특징

- **명시적 의존성 주입**: `main.go`에서 모든 의존성을 수동으로 주입하여 컴포넌트 간 결합 관계를 명확히 파악 가능
- **도메인 응집도**: 관련된 Handler와 Service를 `feature/{domain}` 아래에 모아두어 도메인 수정 시 영향 범위 파악이 용이
- **Type-Safe SQL**: `sqlc`를 도입하여 런타임 쿼리 에러를 방지하고 DB 스키마와 코드의 싱크를 강제함
