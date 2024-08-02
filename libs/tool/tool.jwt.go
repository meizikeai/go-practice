package tool

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mode string = os.Getenv("GO_MODE")

// generate jwt key
// openssl genrsa -out private.key 2048
// openssl rsa -in private.key -pubout -out public.key
var jwtRsaKey map[string]map[string]string = map[string]map[string]string{
	"test": {
		"private": "",
		"public":  "",
	},
	"release": {
		"private": "",
		"public":  "",
	},
}

type JWT struct{}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) EncryptToken(uid, expiration int64) (string, error) {
	if expiration <= 0 {
		expiration = 3196800
	}

	// load the private key
	privateKeyData, _ := Base64DecodeString(jwtRsaKey[mode]["private"])
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)

	times := time.Now().Unix()

	// create a new token object, specifying signing method and claims
	claims := jwt.MapClaims{
		"app": 1,
		"exp": times + expiration,
		"iat": times,
		"uid": HashidsEncode(uid),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// sign the token with the private key
	token, err := t.SignedString(privateKey)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	// fmt.Println("Generated Token:", token)
	return token, nil
}

func (j *JWT) DecryptToken(token string) (map[string]any, error) {
	result := make(map[string]any, 0)

	// load the public key
	publicKeyData, _ := Base64DecodeString(jwtRsaKey[mode]["public"])
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)

	// parse and validate the token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return result, err
	}

	// validate the token
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		result = claims
		fmt.Println("Token is valid! Claims:", string(MarshalJson(result)))
	}

	return result, nil
}
