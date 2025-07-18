package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategorySpendingSummary struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserSessionID uuid.UUID  `gorm:"type:uuid;not null" json:"user_session_id"`
	CategoryID    uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	Category      string     `gorm:"type:varchar(255);not null" json:"category"`
	TotalAmount   int64      `gorm:"type:int8;not null" json:"total_amount"`
	PeriodStart   time.Time  `gorm:"type:timestamp with time zone;not null" json:"period_start"`
	PeriodEnd     time.Time  `gorm:"type:timestamp with time zone;not null" json:"period_end"`
	PeriodType    string     `gorm:"type:varchar(10);not null" json:"period_type"`
	CreatedAt     *time.Time `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt     *time.Time `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at,omitempty"`
}

func (summary *CategorySpendingSummary) BeforeCreate(_ *gorm.DB) error {
	summary.ID = uuid.New()
	return nil
}
