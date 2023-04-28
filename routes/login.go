package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/pkg/auth"
)

func addLoginRoutes(r *gin.Engine, auth *auth.OIDC) {
	r.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusFound, auth.Oauth2Config.AuthCodeURL("random-state"))
	})

	r.GET("/login/callback", func(c *gin.Context) {
		rawIDToken, err := auth.HandleCallback(c.Request.Context(), c.Query("code"), c.Query("state"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": rawIDToken})
	})
}
