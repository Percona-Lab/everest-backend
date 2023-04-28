package routes

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/pkg/auth"
)

func Initialize(ctx context.Context, r *gin.Engine) error {
	oidc, err := auth.New(ctx)
	if err != nil {
		return err
	}

	addLoginRoutes(r, oidc)

	// Everything below needs to be authenticated
	r.Use(ensureAuthenticated(oidc))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": 1})
	})

	return nil
}

func ensureAuthenticated(oidc *auth.OIDC) func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusOK, gin.H{"err": "Invalid authorization header"})
			return
		}

		rawIDToken := auth[len("Bearer "):]
		if _, err := oidc.IDTokenVerifier.Verify(c.Request.Context(), rawIDToken); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}

		c.Next()
	}
}
