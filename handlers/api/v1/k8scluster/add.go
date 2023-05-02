package k8scluster

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"gorm.io/gorm"
)

type K8sCreateRequest struct {
	Name       string `json:"name"`
	ProjectIDs []uint `json:"project_ids"`
}

func Add(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req K8sCreateRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("invalid request. %w", err))
			return
		}

		projects := make([]models.Project, 0, len(req.ProjectIDs))
		for _, p := range req.ProjectIDs {
			projects = append(projects, models.Project{
				Model: gorm.Model{ID: p},
			})
		}

		cluster := models.K8sCluster{
			Name:     req.Name,
			Projects: projects,
		}

		tx := db.Create(&cluster)
		if tx.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot create k8s cluster. %w", tx.Error))
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": cluster.ID})
	}
}
