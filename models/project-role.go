package models

import "gorm.io/gorm"

type ProjectRole struct {
	gorm.Model
	ProjectID int
	Project   Project
	Role      string
	UserID    int
	User      User
}
