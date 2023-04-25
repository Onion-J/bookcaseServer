package api

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
)

type Bookcases struct {
	Area           string `json:"area"`
	AreaName       string `json:"areaName"`
	BookcaseNumber int    `json:"bookcaseNumber"`
}

type BookcaseDeleteList struct {
	DeleteList []Bookcases `json:"deleteList"`
}

// GetBookcaseInfo 获取图书柜情况
func GetBookcaseInfo(c *gin.Context) {
	var bookcaseList []model.Cabinet
	db := common.GetDB()

	// 查询图书柜信息
	if err := db.Find(&bookcaseList).Error; err != nil {
		response.Failed(c, response.NotSelected)
		return
	}

	// 类型转换
	var bookcaseListDto = make([]dto.CabinetDto, len(bookcaseList))
	for i, v := range bookcaseList {
		bookcaseListDto[i] = dto.ToCabinetDto(v)
	}

	// 返回结果
	response.Success(c, response.Selected, gin.H{"bookcaseList": bookcaseListDto})
}

// CreateBookcase 创建图书柜
func CreateBookcase(c *gin.Context) {
	var bookcases Bookcases

	// 参数绑定
	if err := c.ShouldBind(&bookcases); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if bookcases.AreaName == "" || bookcases.BookcaseNumber == 0 {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 判断该区域是否已存在
	if row := db.Where("area = ?", bookcases.AreaName).First(&model.Cabinet{}).RowsAffected; row != 0 {
		response.Failed(c, response.DataAlreadyExists)
		return
	}

	// 创建图书柜
	for i := 1; i <= bookcases.BookcaseNumber; i++ {
		if err := db.Create(&model.Cabinet{
			Area:           bookcases.AreaName,
			SequenceNumber: i,
			Occupied:       false,
		}).Error; err != nil {
			response.Failed(c, response.NotCreated)
			return
		}
	}

	// 返回结果
	response.Success(c, response.Created, nil)
}

// DeleteBookcase 删除图书柜
func DeleteBookcase(c *gin.Context) {

	var bookcaseDeleteList BookcaseDeleteList

	// 参数绑定
	if err := c.ShouldBind(&bookcaseDeleteList); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	for i := 0; i < len(bookcaseDeleteList.DeleteList); i++ {
		if bookcaseDeleteList.DeleteList[i].AreaName == "" {
			response.Failed(c, response.DataError)
			return
		}
	}

	var deleteAreaList = make([]string, len(bookcaseDeleteList.DeleteList))
	for i, bookcase := range bookcaseDeleteList.DeleteList {
		deleteAreaList[i] = bookcase.AreaName
	}

	db := common.GetDB()

	for i := 0; i < len(deleteAreaList); i++ {
		// 判断该区域是否存在
		if row := db.Where("area = ?", deleteAreaList[i]).First(&model.Cabinet{}).RowsAffected; row == 0 {
			response.Failed(c, response.DataDoesNotExist)
			return
		}

		// 删除该区域下的所有图书柜
		if err := db.Where("area = ?", deleteAreaList[i]).Delete(&model.Cabinet{}).Error; err != nil {
			response.Failed(c, response.NotDeleted)
			return
		}
	}

	// 返回结果
	response.Success(c, response.Deleted, nil)
}

// AddBookcase 添加图书柜
func AddBookcase(c *gin.Context) {
	var bookcases Bookcases

	// 参数绑定
	if err := c.ShouldBind(&bookcases); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if bookcases.AreaName == "" || bookcases.BookcaseNumber == 0 {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 判断该区域是否存在
	var bookcaseList []model.Cabinet
	row := db.Where("area = ?", bookcases.AreaName).Find(&bookcaseList).RowsAffected
	if row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 添加图书柜
	for i := 1; i <= bookcases.BookcaseNumber; i++ {
		if err := db.Create(&model.Cabinet{
			Area:           bookcases.AreaName,
			SequenceNumber: i + int(row),
			Occupied:       false,
		}).Error; err != nil {
			response.Failed(c, response.NotCreated)
			return
		}
	}

	// 返回结果
	response.Success(c, response.Created, nil)
}

// ReduceBookcase 删减图书柜
func ReduceBookcase(c *gin.Context) {
	var bookcases Bookcases

	// 参数绑定
	if err := c.ShouldBind(&bookcases); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if bookcases.AreaName == "" || bookcases.BookcaseNumber == 0 {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 判断该区域是否存在
	var bookcaseList []model.Cabinet
	row := db.Where("area = ?", bookcases.AreaName).Find(&bookcaseList).RowsAffected
	if row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 删减图书柜
	for i := 0; i < bookcases.BookcaseNumber; i++ {
		if err := db.Where("area = ? AND sequence_number = ?", bookcases.AreaName, int(row)-i).Delete(&model.Cabinet{}).Error; err != nil {
			response.Failed(c, response.NotDeleted)
			return
		}
	}

	// 返回结果
	response.Success(c, response.Deleted, nil)
}

// RenameBookcase 修改图书柜区域名
func RenameBookcase(c *gin.Context) {
	var bookcases Bookcases

	// 参数绑定
	if err := c.ShouldBind(&bookcases); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	// 判断数据是否为空
	if bookcases.Area == "" || bookcases.AreaName == "" || bookcases.BookcaseNumber == 0 {
		response.Failed(c, response.DataError)
		return
	}

	db := common.GetDB()

	// 判断旧区域名是否存在
	if row := db.Where("area = ?", bookcases.Area).First(&model.Cabinet{}).RowsAffected; row == 0 {
		response.Failed(c, response.DataDoesNotExist)
		return
	}

	// 判断新区域名是否存在
	if row := db.Where("area = ?", bookcases.AreaName).First(&model.Cabinet{}).RowsAffected; row != 0 {
		response.Failed(c, response.DataAlreadyExists)
		return
	}

	// 更新区域名
	if err := db.Model(&model.Cabinet{}).Where("area = ?", bookcases.Area).Update("area", bookcases.AreaName).Error; err != nil {
		response.Failed(c, response.NotUpdated)
		return
	}

	// 返回结果
	response.Success(c, response.Updated, nil)
}
