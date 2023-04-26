package middleware

import (
	"BookcaseServer/common"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
)

func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 auth中间件set的user
		user, _ := c.Get("user")

		db := common.GetDB()

		// 判断用户权限
		var result model.User
		db.Where("teacher_id = ?", user.(model.User).TeacherId).First(&result)
		if result.IsAdmin == false {
			response.Failed(c, response.PermissionDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
