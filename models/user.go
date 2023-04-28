package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Email string `gorm:"uniqueIndex;size:256"`
}

func CreateUserIfNotExists(db *gorm.DB, email string) error {
	user := &User{Email: email}
	tx := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
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
