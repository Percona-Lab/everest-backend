package models

import "gorm.io/gorm"

type DBCluster struct {
	gorm.Model
	Name         string `gorm:"size:256"`
	K8sClusterID int
	K8sCluster   K8sCluster
	ProjectID    int
	Project      Project
}
