package model

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photoUrl" binding:"required"`
	UserId   uint   `json:"user_id"`
}
