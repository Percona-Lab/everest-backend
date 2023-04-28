package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/pkg/auth"
)

func Index(auth *auth.OIDC) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusFound, auth.Oauth2Config.AuthCodeURL("random-state"))
	}
}
