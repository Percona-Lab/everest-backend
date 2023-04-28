package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/routes"
)

func main() {
	ctx := context.TODO()

	r := gin.Default()
	if err := routes.Initialize(ctx, r); err != nil {
		panic(err)
	}

	r.Run(":3000")
}
