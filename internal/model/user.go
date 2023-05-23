package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Printf("[user][HashPassword] generate pwd err: %v", err)
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(pwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	if err != nil {
		log.Printf("[user][HashPassword] comparing hash pwd error: %v", err)
		return err
	}
	return nil
}
