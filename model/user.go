package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `json:"username" gorm:"unique" binding:"required"`
	Email    string  `json:"email" gorm:"unique" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Photos   []Photo `json:"photos" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
