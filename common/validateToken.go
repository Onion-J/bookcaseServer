package common

import (
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
	"strings"
)

func ValidateToken(c *gin.Context) *Claims {
	// 获取authorization header
	tokenString := c.GetHeader("Authorization")

	// 验证 请求头携带的 token 格式 `Bearer token`
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") { // oauth2.0规定Authorization的字符串开头必须要有Bearer
		response.Failed(c, response.PermissionDenied)
		c.Abort() //将这次请求抛弃
		return nil
	}

	tokenString = tokenString[7:]

	// 验证 token 是否有效
	token, claims, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		response.Failed(c, response.InvalidToken)
		c.Abort()
		return nil
	}

	return claims

}
