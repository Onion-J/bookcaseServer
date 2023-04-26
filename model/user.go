package model

type User struct {
	BaseModel
	TeacherId  string `json:"teacherId" gorm:"unique;not null"`
	Name       string `json:"name" gorm:"not null"`
	Password   string `json:"password" gorm:"size:255;not null"`
	Phone      string `json:"phone"`
	IsAdmin    bool   `json:"isAdmin"`
	IsDisabled bool   `json:"isDisabled"`
}
