package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/pkg/roles"
	"gorm.io/gorm"
)

func AddClusterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/cluster/create/:projectID", roles.RequireUserRole(db, roles.ClusterCreate), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/api/cluster/delete/:projectID", roles.RequireUserRole(db, roles.ClusterDelete), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
}
