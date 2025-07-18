package validation

type CreateSpending struct {
	UserSessionID string  `json:"user_session_id" validate:"required,max=50" example:"user_session_id"`
	Category      string  `json:"category" validate:"required,max=50" example:"food"`
	CategoryID    string  `json:"category_id" validate:"required"`
	Name          string  `json:"name" validate:"required,max=50" example:"fake name"`
	Amount        float64 `json:"amount" validate:"required,number,min=0" example:"100.50"`
	Description   string  `json:"description" validate:"omitempty,max=200" example:"fake description"`
	Datetime      string  `json:"datetime" validate:"required" example:"2023-01-01T00:00:00Z"`
	IsConfirm     bool    `json:"is_confirm" validate:"required" example:"true"`
}

// type Spending struct {
// 	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
// 	UserSessionID uuid.UUID  `gorm:"type:uuid;not null" json:"user_session_id"`
// 	Category      string     `gorm:"type:varchar(255);not null" json:"category"`
// 	CategoryID    *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
// 	Name          string     `gorm:"type:varchar(255);not null" json:"name"`
// 	Amount        float64    `gorm:"type:numeric(12,2);not null" json:"amount"`
// 	Description   string     `gorm:"type:text" json:"description,omitempty"`
// 	Datetime      time.Time  `gorm:"type:timestamp with time zone;not null" json:"datetime"`
// 	IsConfirm     bool       `gorm:"type:boolean;default:false;not null" json:"is_confirm"`
// 	CreatedAt     time.Time  `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`
// 	UpdatedAt     time.Time  `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`
// }

type QuerySpending struct {
	Page   int    `validate:"omitempty,number,max=50"`
	Limit  int    `validate:"omitempty,number,max=50"`
	Search string `validate:"omitempty,max=50"`
}

type QuerySpendingSummary struct {
	Page          int    `validate:"omitempty,number,max=50"`
	Limit         int    `validate:"omitempty,number,max=50"`
	Search        string `validate:"omitempty,max=50"`
	UserSessionID string `validate:"required,max=50"`
	PeriodStart   string `validate:"omitempty,datetime"`
	PeriodEnd     string `validate:"omitempty,datetime"`
	PeriodType    string `validate:"omitempty,oneof=daily weekly monthly custom yearly all" example:"daily"`
}
