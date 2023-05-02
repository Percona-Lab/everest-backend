package project

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"gorm.io/gorm"
)

type ProjectCreateRequest struct {
	Name string `json:"name"`
}

func Create(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ProjectCreateRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("invalid request. %w", err))
			return
		}

		project := models.Project{
			Name: req.Name,
		}

		tx := db.Create(&project)
		if tx.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot create project. %w", tx.Error))
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": project.ID})
	}
}
