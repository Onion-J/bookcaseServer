package dto

import "BookcaseServer/model"

type CabinetDto struct {
	Area           string `json:"area"`
	SequenceNumber int    `json:"sequenceNumber"`
	Occupied       bool   `json:"occupied"`
}

func ToCabinetDto(cabinet model.Cabinet) CabinetDto {
	return CabinetDto{
		Area:           cabinet.Area,
		SequenceNumber: cabinet.SequenceNumber,
		Occupied:       cabinet.Occupied,
	}
}
