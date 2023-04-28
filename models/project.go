package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name string `gorm:"size:256"`
}
