package v1alpha1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID   uuid.UUID `gorm:"type:uuid"`
	Name string    `gorm:"unique"`
}

type TenantAPI struct {
	DB *gorm.DB
}

func (t *TenantAPI) AddRoutes(router *gin.Engine) {
	router.GET("/v1alpha1/tenants/", t.list)
	router.GET("/v1alpha1/tenants/:id", t.get)
	router.DELETE("/v1alpha1/tenants/:id", t.delete)
	router.POST("/v1alpha1/tenants/", t.post)
}

func (t *TenantAPI) list(c *gin.Context) {
	tenants := []Tenant{}
	t.DB.Find(&tenants)
	c.IndentedJSON(http.StatusOK, &tenants)
}

func (t *TenantAPI) post(c *gin.Context) {
	tenant := Tenant{}
	err := c.BindJSON(&tenant)
	if err != nil {
		return
	}

	tenant.ID = uuid.New()
	result := t.DB.Create(&tenant)
	handlePostResult(c, result, tenant)
}

func (t *TenantAPI) get(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	tenant := Tenant{ID: pk}

	err := t.DB.First(&tenant).Error

	handleGetResult(c, err, tenant)
}

func (t *TenantAPI) delete(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	err := t.DB.Delete(&Tenant{ID: pk}).Error

	handleDeleteResult(c, err)
}
