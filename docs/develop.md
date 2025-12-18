# 프로젝트 개선 제안 (Development Improvement Plan)

이 문서는 현재 프로젝트의 상태를 분석하고, 향후 유지보수성과 확장성을 위해 개선이 필요한 항목들을 정리합니다.

## 1. 설정 관리 (Configuration Management)

### 🟢 현재 상태 (개선 완료)

- **`cleanenv` 라이브러리 도입**: `internal/config` 패키지에서 `cleanenv`를 사용하여 환경변수와 YAML 파일을 구조체로 매핑합니다.
- **환경 분리**: `dev`, `prod` 환경에 따라 다른 설정 파일(`application-*.yml`)을 로드하며, 환경변수 우선순위를 적용하여 유연하게 동작합니다.
- **CORS 설정 추가**: `application.yml`을 통해 CORS 정책(Allowed Origins, Methods, Headers)을 중앙에서 관리하도록 개선되었습니다.

### 🟡 추가 제안

1. **Secret 관리**: 프로덕션 배포 시 `DB_PASSWORD`, `JWT_SECRET` 등의 민감 정보는 파일이 아닌 환경변수(CI/CD 파이프라인)로 주입하는 원칙을 준수해야 합니다.

## 2. 모듈 의존성 관리 (Module & Dependency)

### 🟢 현재 상태 (개선 완료)

- **`go.mod` 정리**: `cleanenv`, `jwt/v5`, `zap` 등 주요 의존성이 `require` 블록에 명시적으로 선언되어 관리되고 있습니다.

### 🟡 추가 제안

1. **`go mod tidy` 생활화**: 의존성 변경 시마다 `tidy` 명령어를 실행하여 `go.sum` 무결성을 유지합니다.

## 3. 아키텍처 및 코드 구조 (Architecture & Structure)

### 🟢 현재 상태 (개선 완료)

- **고성능 로거(Zap) 도입**: `uber-go/zap`을 도입하여 구조화된 로깅을 구현했습니다.
  - **Dev**: 가독성이 좋은 컬러 콘솔 인코더 사용
  - **Prod**: 시스템 수집에 최적화된 JSON 인코더 사용 (ISO8601 타임스탬프 등)
- **로깅 래퍼(Wrapper) 구현**: `log.MapStr`, `log.MapErr` 등의 헬퍼 함수를 통해 `zap`의 필드 로깅을 쉽고 안전하게 사용할 수 있도록 추상화했습니다.
- **에러 로깅 표준화**: `log.MapErr`를 통해 `zap.Error(err)`를 사용하여 에러를 구조화된 필드로 일관되게 남기고 있습니다.

### 🟡 추가 제안

1. **에러 래핑 (Error Wrapping)**: 에러를 반환(Return)할 때 `fmt.Errorf("...: %w", err)`를 활용하여 호출 스택의 에러 문맥(Context)을 보존하는 컨벤션을 유지해야 합니다.

## 4. 공통 모듈 활용 (Shared Modules)

### 🟡 추가 제안

1. **`internal/shared` 적극 활용**: 여러 도메인에서 반복되는 에러 처리(`errorx`)나 공통 모델(`model`)은 `internal/shared` 패키지를 통해 중앙화하여 중복을 최소화합니다.

## 5. 데이터베이스 마이그레이션 (Database Migrations)

### 🟡 추가 제안

1. **마이그레이션 자동화**: `migrations` 폴더의 스크립트를 배포 파이프라인에서 자동으로 실행하거나, 애플리케이션 시작 시 검증하는 단계를 추가하여 DB 스키마 싱크를 보장해야 합니다.
