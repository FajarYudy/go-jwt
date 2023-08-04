package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Username string `json:"username" gorm:"unique" binding:"required"`
	Email string `json:"email" gorm:"unique" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type UserUpdate struct {
	Name string `json:"name"`	
}

func (user *User) HasPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providerPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providerPassword))
	if err != nil {
		return err
	}
	return nil
}