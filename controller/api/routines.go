package api

import (
	"BookcaseServer/common"
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type File struct {
	Name string `json:"name"`
}

type FileList struct {
	Files []File `json:"files"`
}

// GetSlide 轮播图
func GetSlide(c *gin.Context) {

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

// UploadImg 上传图片
func UploadImg(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		response.Failed(c, response.UploadFailed)
		return
	}

	baseUrl := viper.GetString("slid.baseUrl")

	filename := "." + baseUrl + filepath.Base(file.Filename)

	if err := c.SaveUploadedFile(file, filename); err != nil {
		response.Failed(c, response.UploadFailed)
		return
	}

	response.Success(c, response.UploadSuccess, nil)
}

// DeleteImg 删除图片
func DeleteImg(c *gin.Context) {
	var fileList FileList

	// 参数绑定
	if err := c.ShouldBind(&fileList); err != nil {
		response.Failed(c, response.ParamError)
		return
	}

	baseUrl := viper.GetString("slid.baseUrl")

	for _, file := range fileList.Files {
		fileName := "." + baseUrl + filepath.Base(file.Name)
		if err := os.Remove(fileName); err != nil {
			response.Failed(c, response.NotDeleted)
			return
		}
	}

	response.Success(c, response.Deleted, nil)

}
