package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	dbcluster "github.com/percona/everest-backend/handlers/api/v1/db/cluster"
	k8scluster "github.com/percona/everest-backend/handlers/api/v1/k8s/cluster"
	"github.com/percona/everest-backend/handlers/api/v1/project"
	"github.com/percona/everest-backend/handlers/api/v1/user"
	"github.com/percona/everest-backend/handlers/login"
	"github.com/percona/everest-backend/pkg/auth"
	"github.com/percona/everest-backend/pkg/roles"
	"gorm.io/gorm"
)

func Initialize(ctx context.Context, r *gin.Engine, db *gorm.DB) error {
	oidc, err := auth.New(ctx)
	if err != nil {
		return err
	}

	rolesRegistry := roles.New()

	r.GET("/login", login.Index(oidc))
	r.GET("/login/callback", login.Callback(oidc, db))

	// Everything below needs to be authenticated
	r.Use(ensureAuthenticated(oidc, db))

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("user/add-role", user.AddRole(db))
		apiv1.POST("user/remove-role", user.RemoveRole(db))

		apiv1.POST("project/create", project.Create(db))

		apiv1.GET("k8s/cluster/list", k8scluster.List(db))
		apiv1.POST("k8s/cluster/add", k8scluster.Add(db))

		projectApi := apiv1.Group("project/:projectID")
		{
			projectApi.POST("db/cluster/create",
				requireUserRole(db, roles.ClusterCreate, rolesRegistry),
				dbcluster.Create(db),
			)
			projectApi.POST("db/cluster/delete",
				requireUserRole(db, roles.ClusterDelete, rolesRegistry),
				dbcluster.Delete(db),
			)
			projectApi.GET("db/cluster/list", dbcluster.List(db))
			projectApi.POST("db/cluster/start", dbcluster.Start(db))
			projectApi.POST("db/cluster/stop", dbcluster.Stop(db))
		}
	}

	r.Use(errorHandler())

	return nil
}
