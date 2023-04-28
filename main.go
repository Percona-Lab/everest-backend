package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/percona/everest-backend/models"
	"github.com/percona/everest-backend/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	ctx := context.TODO()

	dsn := "root:my-secret-pw@tcp(127.0.0.1:3306)/everest?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := models.InitializeModels(db); err != nil {
		panic(err)
	}

	r := gin.Default()
	if err := routes.Initialize(ctx, r, db); err != nil {
		panic(err)
	}

	r.Run(":3000")
}
