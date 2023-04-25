package applet

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

const formatTime = "2006-01-02 15:04:05"

// Apply 申请
func Apply(c *gin.Context) {
	var apply model.Apply
	err := c.ShouldBind(&apply)
	if err != nil {
		response.Failed(c, response.ParamError)
		return
	}
	apply.StartDate += " 00:00:00"
	apply.EndDate += " 23:59:59"

	db := common.GetDB()
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
		db.Where("student_id = ?", apply.StudentId).First(&lastApply)
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

// ApplicationRecords 申请记录
func ApplicationRecords(c *gin.Context) {
	studentId := c.PostForm("studentId")
	name := c.PostForm("name")
	db := common.GetDB()
	rows := db.Where("student_id = ? AND name = ?", studentId, name).Find(&model.Apply{}).RowsAffected
	if rows == 0 {
		response.Success(c, response.Selected, gin.H{"records": "none"})
	} else {
		var applyList []model.Apply
		db.Where("student_id = ? AND name = ?", studentId, name).Find(&applyList)
		var records = make([]dto.ApplyDto, len(applyList))
		for index, value := range applyList {
			records[index] = dto.ToApplyDto(value)
		}
		response.Success(c, response.Selected, gin.H{"records": records})
	}

}
