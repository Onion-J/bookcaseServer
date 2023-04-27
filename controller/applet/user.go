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

type StudentData struct {
	Code      string `json:"code"`
	OpenId    string `json:"openId"`
	StudentId string `json:"studentId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var studentData StudentData
	var code2sessionResult model.Code2SessionResult

	// 参数绑定
	if err := c.ShouldBind(&studentData); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if studentData.Code == "" {
		response.Failed(c, response.DataError)
		return
	}

	// 获取openId
	appID := viper.GetString("applet.appID")
	appSecret := viper.GetString("applet.appSecret")
	code2sessionURL := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	url := fmt.Sprintf(code2sessionURL, appID, appSecret, studentData.Code)

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
	var result model.Student
	if rows := db.Where("open_id = ?", code2sessionResult.OpenId).First(&result).RowsAffected; rows == 0 {
		// 不存在
		response.Respond(c, http.StatusAccepted, 202, "请绑定学生账户！", nil)
		return
	}

	// 返回结果
	response.Success(c, response.LoginSuccess, gin.H{"OpenId": code2sessionResult.OpenId, "studentInfo": dto.ToStudentDto(result)})
}

// BindStudentAccount 绑定学生账户
func BindStudentAccount(c *gin.Context) {
	var studentData StudentData

	// 参数绑定
	if err := c.ShouldBind(&studentData); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if studentData.Code == "" || studentData.StudentId == "" || studentData.Name == "" || studentData.Phone == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询studentId是否已经存在
	if rows := db.Where("student_id = ?", studentData.StudentId).First(&model.Student{}).RowsAffected; rows == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	var code2sessionResult model.Code2SessionResult

	// 获取openId
	appID := viper.GetString("applet.appID")
	appSecret := viper.GetString("applet.appSecret")
	code2sessionURL := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	url := fmt.Sprintf(code2sessionURL, appID, appSecret, studentData.Code)

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
	fmt.Println(code2sessionResult.OpenId)

	// 查询openId是否已经存在
	if rows := db.Where("open_id = ?", code2sessionResult.OpenId).First(&model.Student{}).RowsAffected; rows != 0 {
		response.Failed(c, response.DataAlreadyExists)
		return
	}

	if err := db.Model(model.Student{}).Where("student_id = ?", studentData.StudentId).Updates(model.Student{
		OpenId: code2sessionResult.OpenId,
		Name:   studentData.Name,
		Phone:  studentData.Phone,
	}).Error; err != nil {
		response.Failed(c, response.NotUpdated)
		return
	}

	var result model.Student
	if row := db.Where("open_id = ?", code2sessionResult.OpenId).First(&result).RowsAffected; row == 0 {
		response.Failed(c, response.NotSelected)
		return
	}

	// 返回结果
	response.Success(c, response.Updated, gin.H{"OpenId": code2sessionResult.OpenId, "studentInfo": dto.ToStudentDto(result)})
}
