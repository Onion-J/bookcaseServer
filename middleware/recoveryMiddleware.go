package middleware

import (
	"BookcaseServer/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

//  RecoveryMiddleware 拦截panic(err)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Failed(c, fmt.Sprint(err))
			}
		}()

		c.Next()
	}
}
