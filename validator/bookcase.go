package validator

type BookcaseList struct {
	RequestBookcase []CreateBookcaseRequest `json:"requestBookcase" binding:"required"`
}

type CreateBookcaseRequest struct {
	Area           string `json:"area" binding:"required"`
	SequenceNumber int    `json:"sequenceNumber" binding:"required"`
}
