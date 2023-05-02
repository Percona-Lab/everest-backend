package models

import "gorm.io/gorm"

func InitializeModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&DBCluster{},
		&K8sCluster{},
		&Project{},
		&ProjectRole{},
		&User{},
	)
}
