package sec

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// JWT hold secret data it with []byte.
type JWT struct {
	DefExpFunc func() int64
	secret     []byte
}

// NewJWT declare JWT struct with a secret.
func NewJWT(secret []byte, defExpFunc func() int64) *JWT {
	return &JWT{
		secret:     secret,
		DefExpFunc: defExpFunc,
	}
}

// Generate function get custom values and add 'exp' as expires at with expDate argument with unix format.
func (t *JWT) Generate(claims map[string]interface{}, expDate int64) (string, error) {
	mapClaims := jwt.MapClaims{}
	for k := range claims {
		mapClaims[k] = claims[k]
	}

	mapClaims["exp"] = expDate

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		err = fmt.Errorf("cannot sign: %w", err)
	}

	return tokenString, err
}

// Validate is validating and getting claims.
func (t *JWT) Validate(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validate: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token: %w", err)
}
