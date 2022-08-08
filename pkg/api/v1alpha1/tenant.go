package v1alpha1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Tenant struct {
	Base
	Name          string `gorm:"unique" binding:"required"`
	Subscriptions []Subscription
}

type TenantAPI struct {
	DB *gorm.DB
}

func (t *TenantAPI) AddRoutes(router *gin.Engine) {
	router.GET("/v1alpha1/tenants/", t.list)
	router.GET("/v1alpha1/tenants/:tid", t.get)
	router.DELETE("/v1alpha1/tenants/:tid", t.delete)
	router.POST("/v1alpha1/tenants/", t.post)
}

func (t *TenantAPI) list(c *gin.Context) {
	tenants := []Tenant{}
	t.DB.Model(&Tenant{}).Preload("Subscriptions").Find(&tenants)
	c.IndentedJSON(http.StatusOK, &tenants)
}

func (t *TenantAPI) post(c *gin.Context) {
	tenant := Tenant{}
	err := c.BindJSON(&tenant)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result := t.DB.Create(&tenant)
	handlePostResult(c, result, tenant)
}

func (t *TenantAPI) get(c *gin.Context) {
	pk, ok := parseUUID(c, "tid")
	if !ok {
		return
	}

	tenant := Tenant{Base: Base{ID: pk}}

	err := t.DB.Model(&Tenant{}).Preload("Subscriptions").First(&tenant).Error

	handleGetResult(c, err, tenant)
}

func (t *TenantAPI) delete(c *gin.Context) {
	pk, ok := parseUUID(c, "tid")
	if !ok {
		return
	}

	err := t.DB.Delete(&Tenant{Base: Base{ID: pk}}).Error

	handleDeleteResult(c, err)
}
