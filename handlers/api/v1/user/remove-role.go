package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"gorm.io/gorm"
)

type removeRoleRequest struct {
	ProjectID int    `json:"project_id"`
	Role      string `json:"role"`
	UserID    int    `json:"user_id"`
}

func RemoveRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req removeRoleRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("invalid request. %w", err))
			return
		}

		pr := models.ProjectRole{
			ProjectID: req.ProjectID,
			UserID:    req.UserID,
			Role:      req.Role,
		}
		tx := db.Where(&pr).Delete(&pr)
		if tx.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot remove role. %w", tx.Error))
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
