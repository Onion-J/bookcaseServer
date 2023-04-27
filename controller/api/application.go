package api

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
	"time"
)

const formatTime = "2006-01-02 15:04:05"

// GetTodayApplication 获取今天的申请记录
func GetTodayApplication(c *gin.Context) {
	var applyList []model.Apply

	db := common.GetDB()

	min := time.Now().Format("2006-01-02") + " 00:00:00"
	max := time.Now().Format("2006-01-02") + " 23:59:59"

	// 查询
	db.Where("created_at > ? AND created_at < ?", min, max).Find(&applyList)
	// 类型转换
	var applyDtoList = make([]dto.ApplyDto, len(applyList))
	for i, apply := range applyList {
		applyDtoList[i] = dto.ToApplyDto(apply)
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"applyList": applyDtoList})
}

// GetNotExpiredApplication 获取没过期的申请记录
func GetNotExpiredApplication(c *gin.Context) {
	var applyList []model.Apply

	db := common.GetDB()

	// 查询
	db.Where("end_date > ?", time.Now().Format(formatTime)).Find(&applyList)
	// 类型转换
	var applyDtoList = make([]dto.ApplyDto, len(applyList))
	for i, apply := range applyList {
		applyDtoList[i] = dto.ToApplyDto(apply)
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"applyList": applyDtoList})
}

// GetUsageRecords 获取储物柜的使用记录
func GetUsageRecords(c *gin.Context) {
	var cabinet model.Cabinet

	// 参数绑定
	if err := c.ShouldBind(&cabinet); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if cabinet.Area == "" || cabinet.SequenceNumber == 0 {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()
	// 查询
	var result []model.Apply
	db.Where("area = ? AND sequence_number = ?", cabinet.Area, cabinet.SequenceNumber).Find(&result)

	// 类型转换
	var applyDtoList = make([]dto.ApplyDto, len(result))
	for i, apply := range result {
		applyDtoList[i] = dto.ToApplyDto(apply)
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"applyList": applyDtoList})
}

// GetYesterdayNumber 获取一周的申请
func GetYesterdayNumber(c *gin.Context) {
	db := common.GetDB()
	var result []model.Apply

	numberList := make([]int, 7)
	for i := 0; i < 7; i++ {
		min := time.Now().AddDate(0, 0, -(i+1)).Format("2006-01-02") + " 00:00:00"
		max := time.Now().AddDate(0, 0, -(i+1)).Format("2006-01-02") + " 23:59:59"
		row := db.Where("created_at > ? AND created_at < ?", min, max).Find(&result).RowsAffected
		numberList[i] = int(row)
	}

	response.Success(c, response.Selected, gin.H{"numberList": numberList})
}
