package api

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type ChangePasswords struct {
	TeacherId   string `json:"teacherId"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserData struct {
	Users []model.User `json:"users"`
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
		response.Failed(c, "用户不存在！")
		return
	}

	// 检查账户是否被禁用
	if result.IsDisabled {
		response.Failed(c, "账户被禁用，请联系超级管理员！")
		return
	}

	// 检查密码是否一致
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)); err != nil {
		response.Failed(c, "密码错误！")
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

// GetUserInfo 获取登录用户信息
func GetUserInfo(c *gin.Context) {
	// 获取 auth中间件set的user
	user, _ := c.Get("user")

	//返回结果
	response.Success(c, response.VerificationSuccess, gin.H{"user": dto.ToUserDto(user.(model.User))})
}

// GetUserInfoList 获取用户信息
func GetUserInfoList(c *gin.Context) {
	var userList []model.User

	db := common.GetDB()

	// 查询
	if err := db.Find(&userList).Error; err != nil {
		response.Failed(c, response.NotSelected)
		return
	}

	// 类型转换
	var userDtoList = make([]dto.UserDto, len(userList))
	for i, user := range userList {
		userDtoList[i] = dto.ToUserDto(user)
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"userInfoList": userDtoList})
}

// UpdatePhone 更新手机号
func UpdatePhone(c *gin.Context) {
	var user model.User

	// 参数绑定
	if err := c.ShouldBind(&user); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if user.TeacherId == "" || user.Phone == "" {
		response.Failed(c, response.DataError)
		return
	}

	fmt.Println(user)

	db := common.GetDB()

	// 查询用户是否存在以及手机号是否相同
	var result model.User
	if row := db.Where("teacher_id = ?", user.TeacherId).First(&model.User{}).RowsAffected; row == 0 {
		fmt.Println(row)
		response.Failed(c, response.NotUpdated)
		return
	}
	db.Where("teacher_id = ?", user.TeacherId).First(&result)
	if result.Phone == user.Phone {
		fmt.Println(1)
		response.Failed(c, response.NotUpdated)
		return
	}
	fmt.Println(result)
	// 更新手机号
	if err := db.Model(&model.User{}).Where("teacher_id = ?", user.TeacherId).Update("phone", user.Phone).Error; err != nil {
		fmt.Println(2)
		response.Failed(c, response.NotUpdated)
		return
	}

	// 返回结果
	response.Success(c, response.Updated, nil)
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

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var UserData UserData

	// 参数绑定
	if err := c.ShouldBind(&UserData); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if len(UserData.Users) == 0 {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 加密密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	// 创建
	for _, user := range UserData.Users {
		user.Password = string(hashedPassword)
		if err := db.Create(&user).Error; err != nil {
			response.Failed(c, response.NotCreated)
			return
		}
	}

	// 返回结果
	response.Success(c, response.Created, nil)
}

// SetAdmin 设置管理员
func SetAdmin(c *gin.Context) {
	var user model.User

	// 参数绑定
	if err := c.ShouldBind(&user); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if user.TeacherId == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询用户是否存在
	var result model.User
	if row := db.Where("teacher_id = ?", user.TeacherId).First(&result).RowsAffected; row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 判断 isAdmin, true 跳过，false 更新isAdmin为true
	if !result.IsAdmin {
		if err := db.Model(&model.User{}).Where("teacher_id = ?", user.TeacherId).Update("is_admin", true).Error; err != nil {
			response.Failed(c, response.NotUpdated)
			return
		}
	}

	// 返回结果
	response.Success(c, response.Updated, nil)
}

// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {
	var user model.User

	// 参数绑定
	if err := c.ShouldBind(&user); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if user.TeacherId == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询用户是否存在
	if row := db.Where("teacher_id = ?", user.TeacherId).First(&model.User{}).RowsAffected; row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 加密密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	// 更新密码
	if err := db.Model(&model.User{}).Where("teacher_id = ?", user.TeacherId).Update("password", string(hashedPassword)).Error; err != nil {
		response.Failed(c, response.NotUpdated)
		return
	}

	// 返回结果
	response.Success(c, response.Updated, nil)
}

// LockedAccount 禁用账户
func LockedAccount(c *gin.Context) {
	var user model.User

	// 参数绑定
	if err := c.ShouldBind(&user); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if user.TeacherId == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询用户是否存在
	var result model.User
	if row := db.Where("teacher_id = ?", user.TeacherId).First(&result).RowsAffected; row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 判断 IsDisabled, true 跳过，false 更新IsDisabled为true
	if !result.IsDisabled {
		if err := db.Model(&model.User{}).Where("teacher_id = ?", user.TeacherId).Update("is_disabled", true).Error; err != nil {
			response.Failed(c, response.NotUpdated)
			return
		}
	}

	// 返回结果
	response.Success(c, response.Updated, nil)
}

// UnlockedAccount 解禁账户
func UnlockedAccount(c *gin.Context) {
	var user model.User

	// 参数绑定
	if err := c.ShouldBind(&user); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if user.TeacherId == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询用户是否存在
	var result model.User
	if row := db.Where("teacher_id = ?", user.TeacherId).First(&result).RowsAffected; row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 判断 IsDisabled, false 跳过，true 更新IsDisabled为false
	if result.IsDisabled {
		if err := db.Model(&model.User{}).Where("teacher_id = ?", user.TeacherId).Update("is_disabled", false).Error; err != nil {
			response.Failed(c, response.NotUpdated)
			return
		}
	}

	// 返回结果
	response.Success(c, response.Updated, nil)

}

// DeleteAccount 删除账户
func DeleteAccount(c *gin.Context) {
	var user model.User

	// 参数绑定
	if err := c.ShouldBind(&user); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if user.TeacherId == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询用户是否存在
	var result model.User
	if row := db.Where("teacher_id = ?", user.TeacherId).First(&result).RowsAffected; row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 删除
	if err := db.Delete(&result).Error; err != nil {
		response.Failed(c, response.NotDeleted)
		return
	}

	// 返回结果
	response.Success(c, response.Deleted, nil)
}
