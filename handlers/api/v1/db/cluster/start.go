package cluster

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"github.com/percona/everest-backend/pkg/controlplane"
	"gorm.io/gorm"
)

type DBClusterStartRequest struct {
	DBClusterID int `json:"db_cluster_id"`
}

func Start(db *gorm.DB) gin.HandlerFunc {
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

		out, err := controlplane.Start(&cluster)
		if err != nil {
			fmt.Println("**********")
			fmt.Println(out)
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot start cluster. %w", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"out": out})
	}
}
