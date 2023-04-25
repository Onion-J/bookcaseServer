package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Apply struct {
	BasementModel
	StudentId      string `json:"studentId"`
	Name           string `json:"name"`
	Area           string `json:"area"`
	SequenceNumber int    `json:"sequenceNumber"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
}

func (a *Apply) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewV4()
	return
}
