package api

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type ChangePasswords struct {
	TeacherId   string `json:"teacherId"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Login 登录
func Login(c *gin.Context) {
	var user model.User

	// 参数绑定
	if err := c.ShouldBind(&user); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if user.TeacherId == "" || user.Password == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询用户是否存在
	var result model.User
	if row := db.Where("teacher_id = ?", user.TeacherId).First(&result).RowsAffected; row == 0 {
		response.Failed(c, response.LoginFailed)
		return
	}

	// 检查密码是否一致
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)); err != nil {
		response.Failed(c, response.LoginFailed)
		return
	}

	// 发放token 由三段组成，分别为Header表头、Payload负载、Signature签名 ，表头存储元数据（算法名称和令牌类型），负载存放传递的数据，数据分为Public和Private，签名是对Header和Payload两部分的Hash，目的是为了防止数据被篡改
	token, err := common.ReleaseToken(result, common.TokenExpirationTime)
	if err != nil {
		response.Failed(c, response.SystemException)
		log.Printf("token generate error: %v", err)
		return
	}

	// 返回结果
	response.Success(c, response.LoginSuccess, gin.H{"user": dto.ToUserDto(result), "token": token})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 获取 auth中间件set的user
	user, _ := c.Get("user")

	//返回结果
	response.Success(c, response.VerificationSuccess, gin.H{"user": dto.ToUserDto(user.(model.User))})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var passwords ChangePasswords

	// 参数绑定
	if err := c.ShouldBind(&passwords); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if passwords.TeacherId == "" || passwords.OldPassword == "" || passwords.NewPassword == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询用户是否存在
	var result model.User
	if row := db.Where("teacher_id = ?", passwords.TeacherId).First(&result).RowsAffected; row == 0 {
		response.Failed(c, response.NotUpdated)
		return
	}
	// 检查密码是否一致
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(passwords.OldPassword)); err != nil {
		response.Failed(c, response.NotUpdated)
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwords.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Failed(c, response.SystemError)
		return
	}
	passwords.NewPassword = string(hashedPassword)

	// 更新密码
	if err := db.Model(&model.User{}).Where("teacher_id = ?", passwords.TeacherId).Update("password", passwords.NewPassword).Error; err != nil {
		response.Failed(c, response.NotUpdated)
		return
	}

	// 返回结果
	response.Success(c, response.Updated, nil)
}
