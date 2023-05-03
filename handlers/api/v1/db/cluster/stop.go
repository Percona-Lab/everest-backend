package cluster

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"github.com/percona/everest-backend/pkg/controlplane"
	"gorm.io/gorm"
)

type DBClusterStopRequest struct {
	DBClusterID int `json:"db_cluster_id"`
}

func Stop(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DBClusterStartRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("invalid request. %w", err))
			return
		}

		cluster := models.DBCluster{}
		tx := db.Find(&cluster, req.DBClusterID)
		if tx.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot find DB cluster. %w", tx.Error))
			return
		}

		out, err := controlplane.Stop(&cluster)
		if err != nil {
			fmt.Println(out)
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot stop cluster. %w", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"out": out})
	}
}
