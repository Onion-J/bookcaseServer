package model

import (
	uuid "github.com/satori/go.uuid"
)

type Base struct {
	CreatedAt Time  `json:"created_at"`
	UpdatedAt Time  `json:"updated_at"`
	DeletedAt *Time `json:"deleted_at"`
}

type BaseModel struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	CreatedAt Time  `json:"created_at"`
	UpdatedAt Time  `json:"updated_at"`
	DeletedAt *Time `json:"deleted_at"`
}

type BasementModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	CreatedAt Time      `json:"created_at"`
	UpdatedAt Time      `json:"updated_at"`
	DeletedAt *Time     `json:"deleted_at"`
}
