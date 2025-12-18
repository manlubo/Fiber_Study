package util

import "golang.org/x/crypto/bcrypt"

// 문자열 암호화
func HashString(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(str),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// 암호화 검증
func VerifyHashString(str string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}
