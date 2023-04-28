package applet

import (
	"BookcaseServer/common"
	"BookcaseServer/dto"
	"BookcaseServer/model"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Slide 轮播图
func Slide(c *gin.Context) {

	realmName := viper.GetString("server.realmName")
	port := ":" + viper.GetString("server.port")
	baseUrl := viper.GetString("slid.baseUrl")

	fileList := common.FileForEach("." + baseUrl)

	imgUrlList := make([]string, len(fileList))
	for i, file := range fileList {
		imgUrlList[i] = realmName + port + baseUrl + file.Name()
	}

	response.Success(c, response.SendSuccess, gin.H{"imgURL": imgUrlList})
}

// GetBookcaseInfo 获取储物柜信息
func GetBookcaseInfo(c *gin.Context) {
	var bookcaseList []model.Cabinet

	db := common.GetDB()

	db.Where("occupied = ?", false).Find(&bookcaseList)

	var bookcaseDtoList = make([]dto.CabinetDto, len(bookcaseList))
	for i, cabinet := range bookcaseList {
		bookcaseDtoList[i] = dto.ToCabinetDto(cabinet)
	}

	response.Success(c, response.SendSuccess, gin.H{"bookcaseInfo": bookcaseDtoList})
}
