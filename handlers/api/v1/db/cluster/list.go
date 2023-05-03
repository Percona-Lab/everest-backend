package cluster

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"github.com/percona/everest-backend/pkg/controlplane"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListResponse []ListCluster
type ListCluster struct {
	ID         uint           `json:"id"`
	Name       string         `json:"string"`
	K8sCluster ListClusterK8s `json:"k8s_cluster"`
	Status     string         `json:"status"`
}
type ListClusterK8s struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func List(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusters := []models.DBCluster{}
		tx := db.Preload(clause.Associations).Find(&clusters)
		if tx.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot list db clusters. %w", tx.Error))
			return
		}

		res := make(ListResponse, 0, len(clusters))
		for _, c := range clusters {
			status, out, err := controlplane.Status(&c)
			if err != nil {
				status = "unknown"
				fmt.Println(out)
				fmt.Println(err)
			}

			res = append(res, ListCluster{
				ID:   c.ID,
				Name: c.Name,
				K8sCluster: ListClusterK8s{
					ID:   c.K8sCluster.ID,
					Name: c.K8sCluster.Name,
				},
				Status: status,
			})
		}

		c.JSON(http.StatusOK, res)
	}
}
