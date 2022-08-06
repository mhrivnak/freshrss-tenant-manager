package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/saas-patterns/freshrss-tenant-manager/pkg/api"
	"github.com/saas-patterns/freshrss-tenant-manager/pkg/api/v1alpha1"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	for _, model := range []interface{}{
		&v1alpha1.ServiceLevel{},
		&v1alpha1.Tenant{},
	} {
		db.AutoMigrate(model)
	}

	router := gin.Default()

	for _, routeAdder := range []api.RouteAdder{
		&v1alpha1.ServiceLevelAPI{DB: db},
		&v1alpha1.TenantAPI{DB: db},
	} {
		routeAdder.AddRoutes(router)
	}

	router.Run("localhost:8080")
}
