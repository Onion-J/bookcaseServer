package common

import (
	"BookcaseServer/model"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	workDir, _ := os.Getwd() // 获取当前工作目录
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

var DB *gorm.DB

func init() {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.userName")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s", username, password, host, port, database, charset, url.QueryEscape(loc))

	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: driverName,
		DSN:        dsn,
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	if err := db.AutoMigrate(&model.User{}, &model.Student{}, &model.Cabinet{}, &model.Apply{}, &model.Institute{}, &model.Major{}); err != nil {
		panic(err)
		return
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)                  //设置最大连接数
	sqlDB.SetMaxIdleConns(10)                  //设置最大空闲连接数
	sqlDB.SetConnMaxLifetime(90 * time.Second) //空闲连接最多存活时间(超过 90s 建立新的连接,要小于数据库的超时时间)

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

func init() {
	// 创建超级管理员
	var superAdmin model.User
	var password = viper.GetString("superAdmin.password")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	superAdmin = model.User{
		TeacherId:  viper.GetString("superAdmin.teacherId"),
		Name:       viper.GetString("superAdmin.name"),
		Password:   string(hashedPassword),
		IsAdmin:    true,
		IsDisabled: false,
	}

	db := GetDB()
	if row := db.Where("teacher_id = ?", superAdmin.TeacherId).First(&model.User{}).RowsAffected; row == 0 {
		db.Create(&superAdmin)
	}

}
