package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Spending struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserSessionID uuid.UUID  `gorm:"type:uuid;not null" json:"user_session_id"`
	Category      string     `gorm:"type:varchar(255);not null" json:"category"`
	CategoryID    *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	Name          string     `gorm:"type:varchar(255);not null" json:"name"`
	Amount        float64    `gorm:"type:numeric(12,2);not null" json:"amount"`
	Description   string     `gorm:"type:text" json:"description,omitempty"`
	Datetime      time.Time  `gorm:"type:timestamp with time zone;not null" json:"datetime"`
	IsConfirm     bool       `gorm:"type:boolean;default:false;not null" json:"is_confirm"`
	CreatedAt     time.Time  `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`
}

func (spending *Spending) BeforeCreate(_ *gorm.DB) error {
	spending.ID = uuid.New()
	return nil
}
