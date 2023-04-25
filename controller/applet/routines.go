package applet

import (
	"BookcaseServer/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Slid 轮播图
func Slid(c *gin.Context) {
	realmName := viper.GetString("server.realmName")
	port := ":" + viper.GetString("server.port")
	img1 := realmName + port + viper.GetString("slid.img1")
	img2 := realmName + port + viper.GetString("slid.img2")
	img3 := realmName + port + viper.GetString("slid.img3")
	imgURL := []string{img1, img2, img3}
	response.Success(c, response.SendSuccess, gin.H{"imgURL": imgURL})
}
