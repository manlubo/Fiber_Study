# 개발 컨벤션 (Convention)

이 문서는 프로젝트 개발 시 준수해야 할 규칙과 표준을 정의합니다.  
팀원 간의 일관된 코드 품질을 유지하고, 유지보수 효율성을 높이는 것을 목적으로 합니다.

## 1. Naming Convention (네이밍 규칙)

GO 언어의 표준 관례를 따르되, 도메인 명확성을 우선합니다.

### 1.1. 변수 및 함수

- **CamelCase**: 지역 변수, 비공개(unexported) 함수/타입에 사용합니다.
- **PascalCase**: 공개(exported) 함수/타입에 사용합니다.
- **약어 사용 지양**: `req`, `resp`, `ctx` 등 널리 통용되는 약어 외에는 의미가 명확한 단어를 사용합니다.
- **Getter/Setter**: Go 컨벤션에 따라 Getter에 `Get` 접두어를 붙이지 않습니다. (예: `user.Name()` O, `user.GetName()` X)

### 1.2. 파일 및 디렉토리

- **소문자 사용**: 모든 파일명과 디렉토리명은 소문자로 작성합니다.
- **Snake_case**: 복합 명사의 경우 `snake_case`를 사용합니다. (예: `auth_handler.go`, `member_service.go`)
- **테스트 파일**: `_test.go` 접미사를 붙입니다.

### 1.3. 인터페이스

- **-er 접미사**: 메서드가 하나인 경우 `Reader`, `Writer` 등 standard library 관례를 따릅니다.
- **Service/UseCase**: 비즈니스 로직 인터페이스는 명확한 도메인 이름을 사용합니다. (예: `AuthService`)

---

## 2. Directory & Architecture Rules (구조 규칙)

### 2.1. Feature-based Structure

기능 단위로 패키지를 응집시킵니다. `internal/feature/{domain}/` 아래에 관련된 모든 코드를 배치합니다.

- `internal/feature/auth/`
  - `auth_handler.go`: HTTP 핸들러
  - `auth_service.go`: 비즈니스 로직
  - `dto.go`: (권장) Data Transfer Object

### 2.2. 레이어 의존성 규칙 (Strict Layering)

의존성 방향은 항상 **외부 → 내부** 또는 **상위 → 하위** 여야 합니다.

1.  **Handler**: HTTP 요청 처리 외의 로직 금지. **Query(DB)를 직접 호출하지 않습니다.** 반드시 Service를 통합니다.
2.  **Service**: 비즈니스 로직 구현. Handler나 Router를 의존하지 않습니다.
3.  **Query (DataAccess)**: 순수 데이터 처리. 비즈니스 로직을 포함하지 않습니다.

---

## 3. Coding Standards (코딩 표준)

### 3.1. 에러 처리

- **Early Return**: 에러가 발생하면 즉시 반환하여 중첩을 줄입니다.
- **Custom Error**: `study/internal/shared/errorx` 패키지 또는 정의된 상수를 사용하여 에러를 반환합니다.
- **Wrap**: 필요한 경우 context를 추가하여 에러를 래핑합니다.

```go
if err != nil {
    return nil, fmt.Errorf("failed to process: %w", err)
}
```

### 3.2. 로깅 및 Observability

- **Context 사용**: 로그를 남길 때는 반드시 `ctx`를 포함하여 트레이싱 ID가 연동되도록 합니다.
  ```go
  log.InfoCtx(ctx, "회원가입 성공", log.MapStr("email", email))
  ```
- **Span 생성**: 주요 비즈니스 로직 메서드 시작 시 Span을 시작합니다.
  ```go
  ctx, span, start := observability.StartServiceSpan(ctx, "OperationName")
  defer observability.EndSpanWithLatency(span, start, latencyThreshold)
  ```

### 3.3. DTO (Data Transfer Object) 필수 사용

- **규칙**: Handler는 **절대로** DB Model(`query.*`)을 직접 바인딩하지 않습니다. 반드시 API 명세에 맞는 별도의 DTO 구조체를 정의하여 사용합니다.
- **예시**: `SignUpRequest` 등 명확한 의도를 가진 구조체 사용

  ```go
  // Bad
  var req query.Member

  // Good
  var req SignUpRequest
  ```

---

## 4. Git & Commit

- **Commit Message**: `[Feature]`, `[Fix]`, `[Chore]`, `[Docs]` 등의 접두어를 사용하여 커밋의 성격을 명시합니다.
- **Granularity**: 작업 단위별로 쪼개서 커밋합니다. 빌드가 깨지는 커밋은 지양합니다.
