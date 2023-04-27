package router

import (
	"BookcaseServer/controller/api"
	"BookcaseServer/controller/applet"
	"BookcaseServer/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func Start() {
	e := gin.Default()

	// 实现跨域访问
	Cors := cors.New(cors.Config{
		//准许跨域请求网站，多个使用，分开，限制使用*
		AllowOrigins: []string{"*"},
		//准许使用的请求方式
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		//准许使用的请求表头
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
		//显示的请求表头
		ExposeHeaders: []string{"Content-Type"},
		//凭证共享，确定共享
		AllowCredentials: true,
		//容许跨域的原点网站，可以直接return true
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		//超时时间设定
		MaxAge: 24 * time.Hour,
	})

	// 使用跨域和拦截错误中间件
	e.Use(Cors, middleware.RecoveryMiddleware())

	// 设置外部访问静态资源
	e.StaticFS("/resource", http.Dir("./resource"))

	// 小程序路由组
	a := e.Group("/applet")
	// 登录
	a.POST("/login", applet.Login)
	// 绑定学生账户
	a.POST("/bindStudentAccount", applet.BindStudentAccount)
	// 小程序轮播图
	a.GET("/slid", applet.Slid)
	// 获取储物柜信息
	a.GET("/getBookcaseInfo", applet.GetBookcaseInfo)
	// 申请
	a.POST("/apply", applet.Apply)
	// 申请情况
	a.POST("/getApplicationRecords", applet.GetApplicationRecords)
	// 取消申请
	a.POST("/cancelApplication", applet.CancelApplication)

	// Web路由组 "/api"
	w := e.Group("/api")
	// 登录
	w.POST("/login", api.Login)

	// Web路由组>>user路由组 "/api/user"
	u := w.Group("/user")
	// 使用token中间件
	u.Use(middleware.AuthMiddleware())
	// 获取登录用户信息
	u.GET("/getUserInfo", api.GetUserInfo)
	// 获取用户信息
	u.GET("/getUserInfoList", api.GetUserInfoList)
	// 更新手机号
	u.POST("/updatePhone", api.UpdatePhone)
	// 修改密码
	u.POST("/changePassword", api.ChangePassword)
	// 创建用户
	u.POST("/createUser", middleware.PermissionMiddleware(), api.CreateUser)
	// 设置管理员
	u.POST("/setAdmin", middleware.PermissionMiddleware(), api.SetAdmin)
	// 修改用户密码
	u.POST("/resetPassword", middleware.PermissionMiddleware(), api.ResetPassword)
	// 禁用账户
	u.POST("/lockedAccount", middleware.PermissionMiddleware(), api.LockedAccount)
	// 解禁账户
	u.POST("/unlockedAccount", middleware.PermissionMiddleware(), api.UnlockedAccount)
	// 删除账户
	u.POST("/deleteAccount", middleware.PermissionMiddleware(), api.DeleteAccount)

	// Web路由组>>student路由组 "/api/student"
	s := w.Group("/student")
	// 使用token中间件
	//s.Use(middleware.AuthMiddleware())
	// 获取学生入学年份
	s.GET("/getEnrollmentYear", api.GetEnrollmentYear)
	// 创建学生账户
	s.POST("/createStudentAccount", api.CreateStudentAccount)
	// 查询学生账户
	s.POST("/selectStudentAccount", api.SelectStudentAccount)
	// 删除学生账户
	s.POST("/deleteStudentAccount", api.DeleteStudentAccount)

	// Web路由组>>bookcase路由组 "/api/bookcase"
	b := w.Group("/bookcase")
	// 使用token中间件
	b.Use(middleware.AuthMiddleware())
	// 获取储物柜情况
	b.GET("/getBookcaseInfo", api.GetBookcaseInfo)
	// 创建储物柜
	b.POST("/createBookcase", api.CreateBookcase)
	// 删除储物柜
	b.POST("/deleteBookcase", api.DeleteBookcase)
	// 添加储物柜
	b.POST("/addBookcase", api.AddBookcase)
	// 删减储物柜
	b.POST("/reduceBookcase", api.ReduceBookcase)
	// 修改储物柜区域名称
	b.POST("/renameBookcase", api.RenameBookcase)

	// Web路由组>>institute路由组 "/api/institute"
	i := w.Group("/institute")
	// 使用token中间件
	i.Use(middleware.AuthMiddleware())
	// 获取学院及专业情况
	i.GET("/getInstituteInfo", api.GetInstituteInfo)
	// 创建学院及专业
	i.POST("/createInstituteAndMajor", api.CreateInstituteAndMajor)
	// 添加专业
	i.POST("/addMajor", api.AddMajor)
	// 删除专业
	i.POST("/deleteMajor", api.DeleteMajor)
	// 修改学院名称
	i.POST("/renameInstitute", api.RenameInstitute)

	port := viper.GetString("server.port")
	if port != "" {
		panic(e.Run(":" + port))
	}

	panic(e.Run())
}
