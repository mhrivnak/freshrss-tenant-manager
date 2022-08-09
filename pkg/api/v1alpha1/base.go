package v1alpha1

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

type LinkAdder interface {
	AddLinks(base string)
	SelfLink() string
}
