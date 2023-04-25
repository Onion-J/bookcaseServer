package dto

import "BookcaseServer/model"

type StudentDto struct {
	StudentId string `json:"studentId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

func ToStudentDto(user model.Student) StudentDto {
	return StudentDto{
		StudentId: user.StudentId,
		Name:      user.Name,
		Phone:     user.Phone,
	}
}
