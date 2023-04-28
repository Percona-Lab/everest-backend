package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
}
