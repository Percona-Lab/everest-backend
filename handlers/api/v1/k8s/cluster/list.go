package cluster

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListResponse []ListCluster
type ListCluster struct {
	ID       uint                 `json:"id"`
	Name     string               `json:"string"`
	Projects []ListClusterProject `json:"projects"`
}
type ListClusterProject struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func List(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusters := []models.K8sCluster{}
		tx := db.Preload(clause.Associations).Find(&clusters)
		if tx.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot list k8s cluster. %w", tx.Error))
			return
		}

		res := make(ListResponse, 0, len(clusters))
		for _, c := range clusters {
			r := ListCluster{
				ID:       c.ID,
				Name:     c.Name,
				Projects: make([]ListClusterProject, 0, len(c.Projects)),
			}
			fmt.Printf("%#v", r.Projects)
			for _, p := range c.Projects {
				r.Projects = append(r.Projects, ListClusterProject{
					ID:   p.ID,
					Name: p.Name,
				})
			}
			res = append(res, r)
		}

		c.JSON(http.StatusOK, res)
	}
}
