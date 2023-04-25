package model

type Student struct {
	Base
	OpenId    string `json:"openId"`
	StudentId string `json:"studentId" gorm:"primaryKey;type:varchar(12)"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}
