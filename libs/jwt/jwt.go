package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "meizikeai"
)

type Custom struct {
	Uid      int    `json:"uid"`
	UserName string `json:"username"`
}

type customClaims struct {
	Uid      int    `json:"uid"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

const expires = 7 * 24 * time.Hour

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

func GetSignKey() string {
	return SignKey
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

func (j *JWT) EncryptToken(custom Custom) (string, error) {
	claims := customClaims{
		Uid:      custom.Uid,
		UserName: custom.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expires).Unix(),
			Issuer:    "meizikeai@163.com",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) DecryptToken(tokenString string) (*customClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

func (j *JWT) UpdateToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		// claims.StandardClaims.ExpiresAt = time.Now().Add(expires).Unix()

		custom := Custom{
			Uid:      claims.Uid,
			UserName: claims.UserName,
		}

		return j.EncryptToken(custom)
	}

	return "", TokenInvalid
}
