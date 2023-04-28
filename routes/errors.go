package routes

import (
	"log"

	"github.com/gin-gonic/gin"
)

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
