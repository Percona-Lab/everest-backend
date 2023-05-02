package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"index;size:256"`
}

func CreateUserIfNotExists(db *gorm.DB, email string) error {
	user := &User{Email: email}
	tx := db.FirstOrCreate(&user, User{Email: email})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func GetFromEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	tx := db.Where(&User{Email: email}).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}
