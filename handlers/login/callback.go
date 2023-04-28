package login

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"github.com/percona/everest-backend/pkg/auth"
	"gorm.io/gorm"
)

func Callback(auth *auth.OIDC, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawIDToken, claims, err := auth.HandleCallback(c.Request.Context(), c.Query("code"), c.Query("state"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}

		if err = models.CreateUserIfNotExists(db, claims.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": fmt.Errorf("could not save user to db. %w", err),
			})
		}

		c.JSON(http.StatusOK, gin.H{"token": rawIDToken})
	}
}
