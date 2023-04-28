package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"github.com/percona/everest-backend/pkg/auth"
	"github.com/percona/everest-backend/pkg/roles"
	"gorm.io/gorm"
)

func ensureAuthenticated(oidc *auth.OIDC, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		authH := c.GetHeader("Authorization")
		if !strings.HasPrefix(authH, "Bearer ") {
			c.JSON(http.StatusOK, gin.H{"err": "Invalid authorization header"})
			return
		}

		rawIDToken := authH[len("Bearer "):]
		idToken, err := oidc.IDTokenVerifier.Verify(c.Request.Context(), rawIDToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}

		claims := &auth.Claims{}
		if err := idToken.Claims(claims); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": fmt.Errorf("could not parse claims. %w", err)})
			return
		}

		user, err := models.GetFromEmail(db, claims.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": fmt.Errorf("could not find user ID. %w", err)})
			return
		}

		c.Set("userID", user.ID)

		c.Next()
	}
}

func requireUserRole(db *gorm.DB, role roles.Role, registry *roles.Registry) gin.HandlerFunc {
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
			Role:      registry.FindSqlName(role),
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
