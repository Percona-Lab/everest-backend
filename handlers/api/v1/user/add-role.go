package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"github.com/percona/everest-backend/pkg/roles"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type addRoleRequest struct {
	ProjectID int    `json:"project_id"`
	Role      string `json:"role"`
	UserID    int    `json:"user_id"`
}

func AddRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req addRoleRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("invalid request. %w", err))
			return
		}

		_, ok := roles.Get(req.Role)
		if !ok {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("role does not exist"))
			return
		}

		pr := models.ProjectRole{
			ProjectID: req.ProjectID,
			UserID:    req.UserID,
			Role:      req.Role,
		}
		tx := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&pr)
		if tx.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot assign role. %w", tx.Error))
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
