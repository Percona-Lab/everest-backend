package cluster

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"gorm.io/gorm"
)

type DBClusterCreateRequest struct {
	Name         string `json:"name"`
	K8sClusterID int    `json:"k8s_cluster_id"`
}

func Create(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DBClusterCreateRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("invalid request. %w", err))
			return
		}

		projectID, err := strconv.Atoi(c.Param("projectID"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot parse projectID. %w", err))
			return
		}

		cluster := models.DBCluster{
			Name:         req.Name,
			K8sClusterID: req.K8sClusterID,
			ProjectID:    projectID,
		}

		tx := db.Create(&cluster)
		if tx.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot create DB cluster. %w", tx.Error))
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": cluster.ID})
	}
}
