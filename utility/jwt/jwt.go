package jwt

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/golang-jwt/jwt/v5"
	"im/internal/config"
	"time"
)

type MyCustomClaims struct {
	jwt.RegisteredClaims
}

func GenerateJwtToken(id string) string {
	expiresAt := time.Now().Add(time.Second * time.Duration(config.Cfg.Jwt.ExpiresTime))

	myClaims := MyCustomClaims{
		jwt.RegisteredClaims{
			Issuer:    "im.server",
			Subject:   "u-info",
			Audience:  jwt.ClaimStrings{"im"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        id,
		},
	}

	tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims).SignedString([]byte(config.Cfg.Jwt.Secret))
	return tokenString
}

func ParseJwtToken(tokenString string) (jwtClaims *MyCustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, gerror.Newf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Cfg.Jwt.Secret), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return
}
