package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"rental-management/internal/config"
	"rental-management/internal/pkg/errors"
	"rental-management/internal/pkg/response"
)

// Claims JWT claims结构
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint) (string, error) {
	cfg := config.Get()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpireDuration())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.Get()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrTokenInvalid
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, errors.CodeUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.FailWithMsg(c, errors.CodeUnauthorized, "Authorization格式错误")
			c.Abort()
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			if err == errors.ErrTokenExpired {
				response.Fail(c, errors.CodeTokenExpired)
			} else {
				response.Fail(c, errors.CodeTokenInvalid)
			}
			c.Abort()
			return
		}

		// 将用户ID存入上下文
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uint)
	}
	return 0
}
