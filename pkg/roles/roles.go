package roles

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"gorm.io/gorm"
)

type Role int

const (
	ClusterCreate Role = iota
	ClusterDelete
)

var roleRegistry = map[string]Role{
	"cluster.create": ClusterCreate,
	"cluster.delete": ClusterDelete,
}

func Get(role string) (Role, bool) {
	r, ok := roleRegistry[role]
	return r, ok
}

// TODO: fix, this is just bad design
func FindSqlName(role Role) string {
	for k, r := range roleRegistry {
		if r == role {
			return k
		}
	}

	return ""
}

func RequireUserRole(db *gorm.DB, role Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, _ := c.Get("userID")
		userID, ok := v.(uint)
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot assert userID"))
			return
		}

		projectID, err := strconv.Atoi(c.Param("projectID"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot parse projectID. %w", err))
			return
		}

		var pr models.ProjectRole
		tx := db.Where(models.ProjectRole{
			ProjectID: projectID,
			Role:      FindSqlName(role),
			UserID:    int(userID),
		}).First(&pr)

		if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot lookup role. %w", tx.Error))
			return
		}

		if pr.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "User does not have the required role"})
			c.Abort()
			return
		}
	}
}
