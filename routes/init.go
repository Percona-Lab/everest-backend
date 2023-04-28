package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"github.com/percona/everest-backend/pkg/auth"
	"github.com/percona/everest-backend/routes/api"
	"gorm.io/gorm"
)

func Initialize(ctx context.Context, r *gin.Engine, db *gorm.DB) error {
	oidc, err := auth.New(ctx)
	if err != nil {
		return err
	}

	addLoginRoutes(r, oidc, db)

	// Everything below needs to be authenticated
	r.Use(ensureAuthenticated(oidc, db))

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": 1})
	})

	api.AddUserRoutes(r, db)
	api.AddProjectRoutes(r, db)
	api.AddClusterRoutes(r, db)

	r.Use(errorHandler())

	return nil
}

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

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			log.Println(ginErr)
		}

		// status -1 doesn't overwrite existing status code
		c.JSON(-1, gin.H{"err": "see logs"})
	}
}
