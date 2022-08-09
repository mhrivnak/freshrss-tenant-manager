package v1alpha1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Tenant struct {
	Base
	Name          string `gorm:"unique" binding:"required"`
	Subscriptions []Subscription
	Links         TenantLinks `gorm:"-" json:"links"` // ignores this field
}

type TenantLinks struct {
	Self          string `json:"self"`
	Subscriptions string `json:"subscriptions"`
}

func (t *Tenant) AddLinks(base string) {
	t.Links = TenantLinks{
		Self:          fmt.Sprintf("%s/v1alpha1/tenants/%s", base, t.ID.String()),
		Subscriptions: fmt.Sprintf("%s/v1alpha1/tenants/%s/subscriptions/", base, t.ID.String()),
	}

	for i, _ := range t.Subscriptions {
		t.Subscriptions[i].AddLinks(base)
	}
}

func (t *Tenant) SelfLink() string {
	return t.Links.Self
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
	result := t.DB.Model(&Tenant{}).Preload("Subscriptions").Find(&tenants)

	models := make([]LinkAdder, len(tenants))
	for i, _ := range tenants {
		models[i] = &tenants[i]
	}
	handleListResult(c, result.Error, models)
}

func (t *TenantAPI) post(c *gin.Context) {
	tenant := Tenant{}
	err := c.BindJSON(&tenant)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result := t.DB.Create(&tenant)
	handlePostResult(c, result, &tenant)
}

func (t *TenantAPI) get(c *gin.Context) {
	pk, ok := parseUUID(c, "tid")
	if !ok {
		return
	}

	tenant := Tenant{Base: Base{ID: pk}}

	err := t.DB.Model(&Tenant{}).Preload("Subscriptions").First(&tenant).Error

	handleGetResult(c, err, &tenant)
}

func (t *TenantAPI) delete(c *gin.Context) {
	pk, ok := parseUUID(c, "tid")
	if !ok {
		return
	}

	err := t.DB.Delete(&Tenant{Base: Base{ID: pk}}).Error

	handleDeleteResult(c, err)
}
