package dto

import "BookcaseServer/model"

type UserDto struct {
	TeacherId  string `json:"teacherId"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	IsAdmin    bool   `json:"isAdmin"`
	IsDisabled bool   `json:"isDisabled"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		TeacherId:  user.TeacherId,
		Name:       user.Name,
		Phone:      user.Phone,
		IsAdmin:    user.IsAdmin,
		IsDisabled: user.IsDisabled,
	}
}
