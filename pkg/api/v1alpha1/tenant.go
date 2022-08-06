package v1alpha1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	UUID string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type TenantAPI struct {
	DB *gorm.DB
}

func (t *TenantAPI) AddRoutes(router *gin.Engine) {
	router.GET("/v1alpha1/tenants/", t.list)
	router.GET("/v1alpha1/tenants/:uuid", t.get)
	router.DELETE("/v1alpha1/tenants/:uuid", t.delete)
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

	tenant.UUID = uuid.New().String()

	result := t.DB.Create(&tenant)
	if result.Error != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, tenant)
}

func (t *TenantAPI) get(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	tenant := Tenant{}

	err := t.DB.First(&tenant, "uuid = ?", pk).Error

	handleGetResult(c, err, tenant)
}

func (t *TenantAPI) delete(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	err := t.DB.Delete(&Tenant{UUID: pk}).Error

	handleDeleteResult(c, err)
}
