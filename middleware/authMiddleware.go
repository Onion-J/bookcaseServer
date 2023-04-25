package middleware

import (
	"BookcaseServer/common"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := common.ValidateToken(c)
		if claims == nil {
			return
		}

		// 通过验证后获取claim 中的userId
		userId := claims.UserId
		db := common.GetDB()
		var user model.User
		db.First(&user, userId)

		// 用户不存在
		if user.ID == 0 {
			response.Failed(c, response.PermissionDenied)
			c.Abort()
			return
		}

		//用户存在 将user 的信息写入上下文
		c.Set("user", user)
		c.Next()
	}

}
