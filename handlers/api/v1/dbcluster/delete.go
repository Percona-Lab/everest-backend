package dbcluster

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"gorm.io/gorm"
)

type DBClusterDeleteRequest struct {
	DBClusterID int `json:"db_cluster_id"`
}

func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DBClusterDeleteRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("invalid request. %w", err))
			return
		}

		cluster := &models.DBCluster{}
		tx := db.Delete(&cluster, req.DBClusterID)
		if tx.Error != nil {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("cannot delete DB cluster. %w", tx.Error))
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
