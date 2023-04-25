package applet

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

// AppletLogin 用户登录
func AppletLogin(c *gin.Context) {

	var student model.Student
	var code2sessionResult model.Code2SessionResult

	// 获取code
	code := c.PostForm("code")
	appID := viper.GetString("applet.appID")
	appSecret := viper.GetString("applet.appSecret")
	code2sessionURL := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	url := fmt.Sprintf(code2sessionURL, appID, appSecret, code)

	// 发送凭证检验请求
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		fmt.Println("微信登录凭证校验接口请求错误!")
		response.Failed(c, response.LoginFailed)
		return
	}

	// 将返回的数据解析
	if err := json.NewDecoder(res.Body).Decode(&code2sessionResult); err != nil {
		fmt.Println("参数解析错误!")
		response.Failed(c, response.LoginFailed)
		return
	}

	// 查询用户是否已经存在
	db := common.GetDB()
	rows := db.Where("open_id = ?", code2sessionResult.OpenId).First(&model.Student{}).RowsAffected
	if rows == 0 {
		// 不存在，添加用户
		student.OpenId = code2sessionResult.OpenId
		row := db.Create(&student).RowsAffected
		if row == 0 {
			fmt.Println("用户添加失败!")
			response.Failed(c, response.LoginFailed)
			return
		}
	}

	response.Success(c, response.LoginSuccess, gin.H{"OpenId": code2sessionResult.OpenId})
}

// VerifyUser 校验用户
func VerifyUser(c *gin.Context) {
	openId := c.PostForm("openId")
	db := common.GetDB()
	var student model.Student
	db.Where("open_id = ?", openId).Find(&student)
	if student.StudentId == "" || student.Name == "" || student.Phone == "" {
		response.Failed(c, response.VerificationFailed)
		return
	}
	response.Success(c, response.VerificationSuccess, gin.H{"User": dto.ToStudentDto(student)})
}

// GetUserInfo 完善用户信息
func GetUserInfo(c *gin.Context) {
	openId := c.PostForm("openId")
	name := c.PostForm("name")
	studentId := c.PostForm("studentId")
	phone := c.PostForm("phone")
	db := common.GetDB()
	db.Model(&model.Student{}).Where("open_id = ?", openId).Updates(model.Student{
		StudentId: studentId,
		Name:      name,
		Phone:     phone,
	})
	var student model.Student
	db.Where("open_id = ?", openId).Find(&student)
	if student.StudentId == "" || student.Name == "" || student.Phone == "" {
		response.Failed(c, response.NotUpdated)
		return
	}
	response.Success(c, response.Updated, gin.H{"Student": student})
}
