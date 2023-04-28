package models

import "gorm.io/gorm"

func InitializeModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&Project{},
		&ProjectRole{},
		&User{},
	)
}
