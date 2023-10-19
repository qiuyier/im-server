package jwt

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
	"im/internal/config"
	"im/internal/consts"
	"time"
)

type MyCustomClaims struct {
	JwtKeyName string
	jwt.RegisteredClaims
}

func GenerateJwtToken(id, jwtKeyName string) string {
	// 先从redis获取是否存在
	cacheKey := fmt.Sprintf(consts.JwtRedisCacheKeyFormat, jwtKeyName, id)
	res, err := g.Redis().Get(gctx.New(), cacheKey)
	if err != nil || res.IsNil() {
		expiresAt := time.Now().Add(time.Second * time.Duration(config.Cfg.Jwt.ExpiresTime))

		myClaims := MyCustomClaims{
			JwtKeyName: jwtKeyName,
			RegisteredClaims: jwt.RegisteredClaims{
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

		// 使用redis缓存
		_ = grpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
			_, _ = g.Redis().Set(ctx, cacheKey, tokenString)

			_, _ = g.Redis().Expire(ctx, cacheKey, gconv.Int64(config.Cfg.Jwt.ExpiresTime))

		}, nil)
		return tokenString
	}

	return res.String()
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
