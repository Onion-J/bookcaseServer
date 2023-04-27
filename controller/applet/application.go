package applet

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"time"
)

const formatTime = "2006-01-02 15:04:05"

// Apply 申请
func Apply(c *gin.Context) {
	var apply model.Apply

	// 参数绑定
	if err := c.ShouldBind(&apply); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	apply.StartDate += " 00:00:00"
	apply.EndDate += " 23:59:59"

	db := common.GetDB()

	// 判断储物柜是否被占用
	var result model.Cabinet

	db.Where("area = ? AND sequence_number = ?", apply.Area, apply.SequenceNumber).First(&result)
	if result.Occupied {
		response.Failed(c, response.ApplyFailed)
		return
	}

	rows := db.Where("student_id = ?", apply.StudentId).First(&model.Apply{}).RowsAffected
	//  没有记录
	if rows == 0 {
		row := db.Create(&apply).RowsAffected
		if row == 0 {
			response.Failed(c, response.NotCreated)
			return
		}
	} else {
		// 有记录，判断是否还在使用期
		var lastApply model.Apply
		db.Where("student_id = ?", apply.StudentId).Last(&lastApply)
		fmt.Println(lastApply)
		lastEndDate, _ := time.ParseInLocation(formatTime, lastApply.EndDate, time.Local)
		fmt.Println(lastEndDate)

		fmt.Println(lastEndDate.After(time.Now()))
		if lastEndDate.After(time.Now()) {
			response.Failed(c, response.ApplyFailed)
			return
		} else {
			row := db.Create(&apply).RowsAffected
			if row == 0 {
				response.Failed(c, response.NotCreated)
				return
			}
		}
	}
	db.Model(&model.Cabinet{}).Where("area = ? and sequence_Number = ?", apply.Area, apply.SequenceNumber).Update("occupied", true)
	response.Success(c, response.ApplySuccess, nil)
}

// GetApplicationRecords 申请记录
func GetApplicationRecords(c *gin.Context) {
	var student model.Student

	// 参数绑定
	if err := c.ShouldBind(&student); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if student.StudentId == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()
	// 查询
	var result model.Student
	if err := db.Preload(clause.Associations).Where("student_id = ?", student.StudentId).Find(&result).Error; err != nil {
		response.Failed(c, response.NotSelected)
		return
	}

	// 类型转换
	var records = make([]dto.ApplyDto, len(result.Apply))
	for i, apply := range result.Apply {
		records[i] = dto.ToApplyDto(apply)
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"records": records})

}

// CancelApplication 取消申请
func CancelApplication(c *gin.Context) {
	var apply model.Apply

	// 参数绑定
	if err := c.ShouldBind(&apply); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if apply.ID.String() == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()
	// 查询
	var result model.Apply
	if row := db.Where("id = ?", apply.ID).Last(&result).RowsAffected; row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	//更新
	if err := db.Model(&result).Update("end_date", time.Now().Format(formatTime)).Error; err != nil {
		response.Failed(c, response.NotUpdated)
		return
	}

	if err := db.Model(&model.Cabinet{}).Where("area = ? AND sequence_number = ?", result.Area, result.SequenceNumber).Update("occupied", false).Error; err != nil {
		response.Failed(c, response.NotUpdated)
		return
	}

	// 返回结果
	response.Success(c, response.Updated, nil)
}
