package response

import (
	"time"

	"github.com/google/uuid"
)

type SummarySpending struct {
	CategoryID  uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	Category    string     `gorm:"type:varchar(255);not null" json:"category"`
	TotalAmount int64      `gorm:"type:int8;not null" json:"total_amount"`
	PeriodStart time.Time  `gorm:"type:timestamp with time zone;not null" json:"period_start"`
	PeriodEnd   time.Time  `gorm:"type:timestamp with time zone;not null" json:"period_end"`
	PeriodType  string     `gorm:"type:varchar(10);not null" json:"period_type"`
	CreatedAt   *time.Time `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at,omitempty"`
}

type TotalSummarySpending struct {
	Total int64 `json:"total"`
}

type CreateSpending struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Name       string     `json:"name" validate:"required"`
	Amount     int64      `json:"amount" validate:"required"`
	CategoryID uuid.UUID  `json:"category_id" validate:"required"`
	Category   string     `json:"category" validate:"required"`
	Date       time.Time  `json:"datetime" validate:"required"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}
