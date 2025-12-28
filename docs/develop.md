# 개발 및 개선 가이드

이 문서는 프로젝트의 현재 구조적 한계와 이를 개선하기 위한 기술적 제안을 정리합니다.  
현재 프로젝트는 빠르게 기능을 구현하는 데 최적화되어 있으나, 장기적인 유지보수와 안정성을 위해 개선해야 할 지점들이 존재합니다.

## 1. 현재 구조의 한계 및 문제점

### 1.1. Handler와 Data Access Layer의 강한 결합 (개선 중)

- **현상**: 일부 핸들러에서 DB 엔티티(`query.Member` 등)를 직접 요청 바인딩에 사용해왔으나, **최근 `Auth` 도메인 등에서 DTO 도입을 시작했습니다.** 여전히 레거시 코드나 신규 기능 구현 시 주의가 필요합니다.
- **문제점**:
  - **보안 리스크**: DB 스키마 노출 위험 (DTO 미사용 시)
  - **유지보수**: DB 컬럼 변경이 API 명세에 영향을 줌
  - **검증 분리 미흡**: DB 제약조건과 API 유효성 검증의 혼재

### 1.2. Service 계층의 추상화 부족 (Leaky Abstraction)

- **현상**: Service 메서드의 시그니처가 `sqlc` 생성 타입(`query.*`)이나 `pgtype.*`에 의존하고 있습니다.
  ```go
  func (s *AuthService) Register(ctx context.Context, m *query.Member) error
  ```
- **문제점**:
  - Service 계층이 특정 Data Access 라이브러리(`sqlc`, `pgx`)에 종속됩니다.
  - 추후 다른 저장소로 변경하거나 로직을 테스트할 때, `pgtype`과 같은 DB 전용 타입을 계속 다뤄야 합니다.

### 1.3. Service와 Query의 모호한 경계

- **현상**: 비즈니스 로직과 데이터 접근 로직이 섞여 있는 경우가 있습니다. 단순 조회 로직이 Service를 거치면서 불필요한 레이어를 탈 수도 있습니다.

---

## 2. 개선 방향 제안

### 2.1. [필수] DTO (Data Transfer Object) 표준화

API 계층과 데이터 계층의 분리를 위해 **DTO 사용을 전면 표준화**해야 합니다. 현재 `Auth` 도메인에 적용된 `SignUpRequest` 패턴을 다른 도메인으로 확장합니다.

- **흐름**: `Handler` → `Request DTO` → `Service` → `Domain/Response DTO`
- **구체적인 액션**:
  - 신규 기능 개발 시 반드시 DTO 정의
  - 기존 핸들러 리팩토링 시 `query.*` 구조체 사용 제거

### 2.2. [중기] Domain Model 정의

`sqlc`가 생성해주는 `query.Member`는 본질적으로 DB 테이블 스키마입니다. 비즈니스 로직의 주체가 되는 **순수 도메인 모델**을 정의하는 것을 권장합니다.

- **위치**: `internal/domain/member.go` 또는 `internal/feature/member/model.go`
- **내용**: `pgtype` 등이 없는 순수 Go 구조체
- Service는 이 도메인 모델을 입력받고 반환하도록 수정하여 Infra(`pgx`, `sqlc`) 의존성을 제거합니다.

### 2.3. [중기] Service Interface 및 Mocking 강화

현재 Service는 구체 타입(`*AuthService`)으로 주입되고 있습니다. 이를 인터페이스로 추상화하면 테스트가 용이해집니다.

- **방향**:
  - `AuthUseCase` 인터페이스 정의
  - Handler는 `AuthUseCase` 인터페이스에 의존
  - 이를 통해 DB 연결 없이 Handler 단위 테스트 가능

---

## 3. 요약 및 우선순위

| 우선순위 | 구분     | 내용                         | 기대효과                                             |
| :------: | -------- | ---------------------------- | ---------------------------------------------------- |
|  **P1**  | **구조** | **전 도메인 DTO 표준화**     | 보안 강화, API 버저닝 용이성 확보                    |
|    P2    | 테스트   | Service Layer Unit Test 작성 | 비즈니스 로직 안정성 확보                            |
|    P3    | 구조     | Domain Model 분리            | Infrastructure 종속성 제거 (Clean Architecture 지향) |

`Auth` 도메인에서 시작된 **DTO 도입(P1)** 을 전체 프로젝트 컨벤션으로 정착시키는 것이 최우선 과제입니다.
