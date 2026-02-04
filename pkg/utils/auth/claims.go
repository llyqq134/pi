package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const JWTSecret = "secret123"
const JWTExpirationTime = 24 * time.Hour

type Claims struct {
	WorkerID       string `json:"worker_id"`
	Name           string `json:"name"`
	JobTitle       string `json:"jobtitle"`
	DepartmentID   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
	AccessLevel    int    `json:"accesslevel"`
	jwt.RegisteredClaims
}

func GenerateToken(workerID, name, jobTitle, departmentID, departmentName string, accessLevel int) (string, error) {
	expirationTime := time.Now().Add(JWTExpirationTime)

	claims := &Claims{
		WorkerID:       workerID,
		Name:           name,
		JobTitle:       jobTitle,
		DepartmentID:   departmentID,
		DepartmentName: departmentName,
		AccessLevel:    accessLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
