package model

type Institute struct {
	ID    string  `json:"id" gorm:"primaryKey;type:varchar(2)"`
	Name  string  `json:"name" gorm:"unique;not null"`
	Major []Major `json:"major"`
}

type Major struct {
	InstituteID string `json:"instituteId" gorm:"primaryKey"`
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"unique;not null"`
}
