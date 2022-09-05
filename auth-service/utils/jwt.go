package utils

import (
	"errors"
	models "github.com/sukenda/golang-krakend/auth-service/model"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtWrapper struct {
	Kid             string
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type jwtClaims struct {
	jwt.StandardClaims
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

func (w *JwtWrapper) GenerateToken(user models.User) (accessToken, refreshToken string, exp int64, err error) {
	valid := time.Now().Local().Add(time.Hour * time.Duration(w.ExpirationHours)).Unix()
	claims := &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: valid,
			Issuer:    w.Issuer,
		},
		Id:    user.Id,
		Email: user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = w.Kid // This used for krakend and must same with kid on url /auth/jwk

	accessToken, err = token.SignedString([]byte(w.SecretKey))
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, accessToken, valid, err
}

func (w *JwtWrapper) ValidateToken(signedToken string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(w.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)

	if !ok {
		return nil, errors.New("Couldn't parse claims ")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil

}
