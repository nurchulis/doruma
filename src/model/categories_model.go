package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	IsActive  bool      `gorm:"type:boolean;default:true;not null" json:"is_active"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`
}

func (category *Category) BeforeCreate(_ *gorm.DB) error {
	category.ID = uuid.New()
	return nil
}
