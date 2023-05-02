package models

import "gorm.io/gorm"

type K8sCluster struct {
	gorm.Model
	Name     string    `gorm:"size:256"`
	Projects []Project `gorm:"many2many:k8scluster_project"`
}
