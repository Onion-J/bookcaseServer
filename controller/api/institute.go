package api

import (
	"BookcaseServer/common"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type MajorAddList struct {
	AddList []model.Major `json:"addList"`
}

// GetInstituteInfo 获取学院及专业情况
func GetInstituteInfo(c *gin.Context) {
	var institute []model.Institute
	db := common.GetDB()

	// 查询学院及专业信息
	if err := db.Preload(clause.Associations).Find(&institute).Error; err != nil {
		response.Failed(c, response.NotSelected)
		return
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"instituteInfo": institute})
}

// CreateInstituteAndMajor 创建学院及专业
func CreateInstituteAndMajor(c *gin.Context) {
	var institute model.Institute

	// 参数绑定
	if err := c.ShouldBind(&institute); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if institute.ID == "" || institute.Name == "" {
		response.Failed(c, response.DataError)
		return
	}
	if institute.Major != nil {
		for i := 0; i < len(institute.Major); i++ {
			if institute.Major[i].ID == "" || institute.Major[i].Name == "" {
				response.Failed(c, response.DataError)
				return
			}
		}
	}

	db := common.GetDB()

	// 查询学院编号或名称是否存在
	if row := db.Where("id = ? OR name = ?", institute.ID, institute.Name).First(&model.Institute{}).RowsAffected; row != 0 {
		response.Failed(c, response.DataAlreadyExists)
		return
	}

	// 查询专业名称是否存在
	for i := 0; i < len(institute.Major); i++ {
		if row := db.Where("name = ?", institute.Major[i].Name).First(&model.Major{}).RowsAffected; row != 0 {
			response.Failed(c, response.DataAlreadyExists)
			return
		}
	}

	// 创建学院及专业
	if err := db.Create(&institute).Error; err != nil {
		response.Failed(c, response.NotCreated)
		return
	}

	// 返回结果
	response.Success(c, response.Created, nil)
}

// AddMajor 添加专业
func AddMajor(c *gin.Context) {
	var majorAddList MajorAddList

	// 参数绑定
	if err := c.ShouldBind(&majorAddList); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	for i := 0; i < len(majorAddList.AddList); i++ {
		if majorAddList.AddList[i].InstituteID == "" || majorAddList.AddList[i].ID == "" || majorAddList.AddList[i].Name == "" {
			response.Failed(c, response.DataError)
			return
		}
	}

	db := common.GetDB()

	// 判断专业是否存在
	for i := 0; i < len(majorAddList.AddList); i++ {
		if row := db.Where("institute_id = ? AND id = ?", majorAddList.AddList[i].InstituteID, majorAddList.AddList[i].ID).First(&model.Major{}).RowsAffected; row != 0 {
			response.Failed(c, response.DataAlreadyExists)
			return
		}
		if row := db.Where("name = ?", majorAddList.AddList[i].Name).First(&model.Major{}).RowsAffected; row != 0 {
			response.Failed(c, response.DataAlreadyExists)
			return
		}
	}

	// 添加专业
	if err := db.Create(&majorAddList.AddList).Error; err != nil {
		response.Failed(c, response.NotCreated)
		return
	}

	// 返回结果
	response.Success(c, response.Created, nil)
}

// DeleteMajor 删除专业
func DeleteMajor(c *gin.Context) {
	var major model.Major

	// 参数绑定
	if err := c.ShouldBind(&major); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if major.InstituteID == "" || major.ID == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询该专业是否存在
	if row := db.Where("institute_id = ? AND id = ?", major.InstituteID, major.ID).First(&model.Major{}).RowsAffected; row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 删除该专业
	if err := db.Where("institute_id = ? AND id = ?", major.InstituteID, major.ID).Delete(&model.Major{}).Error; err != nil {
		response.Failed(c, response.NotDeleted)
		return
	}

	// 返回结果
	response.Success(c, response.Deleted, nil)

}

// RenameInstitute 修改学院名称
func RenameInstitute(c *gin.Context) {
	var institute model.Institute

	// 参数绑定
	if err := c.ShouldBind(&institute); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if institute.ID == "" || institute.Name == "" {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 查询该学院是否存在
	row1 := db.Where("id = ?", institute.ID).First(&model.Institute{}).RowsAffected
	row2 := db.Where("name = ?", institute.Name).First(&model.Institute{}).RowsAffected
	if row1 == 0 || row2 != 0 {
		response.Failed(c, response.DataError)
		return
	}

	// 更新学院名
	if err := db.Model(&model.Institute{}).Where("id = ?", institute.ID).Update("name", institute.Name).Error; err != nil {
		response.Failed(c, response.NotUpdated)
		return
	}

	// 返回结果
	response.Success(c, response.Updated, nil)
}
