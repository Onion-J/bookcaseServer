package dto

import (
	"BookcaseServer/model"
	uuid "github.com/satori/go.uuid"
)

type ApplyDto struct {
	ID             uuid.UUID `json:"id"`
	StudentId      string    `json:"studentId"`
	Name           string    `json:"name"`
	Area           string    `json:"area"`
	SequenceNumber int       `json:"sequenceNumber"`
	StartDate      string    `json:"startDate"`
	EndDate        string    `json:"endDate"`
}

func ToApplyDto(apply model.Apply) ApplyDto {
	return ApplyDto{
		ID:             apply.ID,
		StudentId:      apply.StudentId,
		Name:           apply.Name,
		Area:           apply.Area,
		SequenceNumber: apply.SequenceNumber,
		StartDate:      apply.StartDate,
		EndDate:        apply.EndDate,
	}
}
