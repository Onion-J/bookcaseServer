package model

type Cabinet struct {
	Base
	Area           string  `json:"area" gorm:"primaryKey;type:varchar(64);not null"`
	SequenceNumber int     `json:"sequenceNumber" gorm:"primaryKey;autoIncrement:false;not null"`
	Occupied       bool    `json:"occupied"`
	Apply          []Apply `json:"apply" gorm:"foreignKey:Area,SequenceNumber;references:Area,SequenceNumber"`
}
