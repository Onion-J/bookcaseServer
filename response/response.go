package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ParamInvalid = "请求参数无效!"
	ParamError   = "参数错误!"

	UploadSuccess = "上传成功!"
	UploadFailed  = "上传失败!"

	ApplySuccess = "申请成功!"
	ApplyFailed  = "申请失败!"

	DataError         = "数据有误!"
	DataAlreadyExists = "数据已存在!"
	DataDoesNotExist  = "数据不存在!"

	Created  = "已创建!"
	Deleted  = "已删除!"
	Updated  = "已更新!"
	Selected = "查询成功!"

	NotCreated  = "创建失败!"
	NotDeleted  = "删除失败!"
	NotUpdated  = "更新失败!"
	NotSelected = "查询失败!"

	SendSuccess = "发送成功!"
	SendFailed  = "发送失败!"

	LoginSuccess = "登录成功!"
	LoginFailed  = "登录失败!"

	SaveSuccess = "保存成功!"
	SaveFailed  = "保存失败!"

	VerificationSuccess = "校验成功!"
	VerificationFailed  = "校验失败!"

	PermissionDenied = "权限不足!"
	InvalidToken     = "token无效!"
	SystemException  = "系统异常!"
	SystemError      = "系统错误!"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

// Success 请求成功返回
func Success(ctx *gin.Context, message string, data gin.H) {
	ctx.JSON(http.StatusOK, Response{200, message, data})
}

// Failed 请求失败返回
func Failed(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{400, message, nil})
}

// Respond 返回请求状态
func Respond(ctx *gin.Context, httpStatus int, code int, message string, data gin.H) {
	ctx.JSON(httpStatus, Response{code, message, data})
}
