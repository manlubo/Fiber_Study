package errorx

import "errors"

// 서비스 에러
var (
	// JSON 파싱 실패
	ErrRequestParseFailed = errors.New("REQUEST_PARSE_FAILED")

	// 필수값 누락
	ErrRequiredFieldMissing = errors.New("REQUIRED_FIELD_MISSING")
)
