package dto

import "BookcaseServer/model"

type UserDto struct {
	TeacherId string `json:"teacherId"`
	Name      string `json:"name"`
	IsAdmin   bool   `json:"isAdmin"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		TeacherId: user.TeacherId,
		Name:      user.Name,
		IsAdmin:   user.IsAdmin,
	}
}
