package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SpendingService interface {
	CreateSpending(c *fiber.Ctx, req *validation.CreateSpending) (*model.Spending, error)
	GetCategories(c *fiber.Ctx, params *validation.QueryUser) ([]model.Category, int64, error)
	GetSpendings(c *fiber.Ctx, params *validation.QueryUser) ([]model.Spending, int64, error)
	GetSummarySpending(c *fiber.Ctx, params *validation.QuerySpendingSummary) ([]model.CategorySpendingSummary, int64, error)
}

type spendingService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewSpendingService(db *gorm.DB, validate *validator.Validate) SpendingService {
	return &spendingService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *spendingService) CreateSpending(c *fiber.Ctx, req *validation.CreateSpending) (*model.Spending, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	userSessionUUID, err := utils.ParseUUID(req.UserSessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid UserSessionID format")
	}

	// Ensure category exists or create it if not
	category := model.Category{Name: req.Category}
	if err := s.DB.WithContext(c.Context()).FirstOrCreate(&category, model.Category{Name: req.Category}).Error; err != nil {
		s.Log.Errorf("Failed to get or create category: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get or create category")
	}
	// Use current time for spending.Datetime
	parsedDatetime := time.Now()

	spending := &model.Spending{
		UserSessionID: userSessionUUID, // userSessionUUID should be uuid.UUID type
		Amount:        req.Amount,
		Name:          req.Name,
		Description:   req.Description,
		Category:      req.Category,
		Datetime:      parsedDatetime,
		CategoryID:    &category.ID, // Pass pointer to uuid.UUID
	}

	result := s.DB.WithContext(c.Context()).Create(spending)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Spending record already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create spending: %+v", result.Error)
	}

	// Call UpsertSummary to update weekly and monthly summaries in background
	go func(db *gorm.DB, userSessionID uuid.UUID, categoryID uuid.UUID, categoryName string, amount int64, log *logrus.Logger, ctx fiber.Ctx) {
		// Use the same context for DB
		err := UpsertSummary(db.WithContext(ctx.Context()), userSessionID, categoryID, categoryName, amount)
		if err != nil {
			log.Errorf("Failed to upsert summary: %+v", err)
			// Tidak perlu return error jika summary gagal — opsional
		}
	}(s.DB, userSessionUUID, category.ID, category.Name, int64(req.Amount), s.Log, *c)

	return spending, result.Error
}

func (s *spendingService) GetSpendings(c *fiber.Ctx, params *validation.QueryUser) ([]model.Spending, int64, error) {
	var spendings []model.Spending
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Model(&model.Spending{}).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR category LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	result := query.Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to count spendings: %+v", result.Error)
		return nil, 0, result.Error
	}

	result = query.Limit(params.Limit).Offset(offset).Find(&spendings)
	if result.Error != nil {
		s.Log.Errorf("Failed to get spendings: %+v", result.Error)
		return nil, 0, result.Error
	}

	return spendings, totalResults, nil
}

// GetCategories implements SpendingService.
func (s *spendingService) GetCategories(c *fiber.Ctx, params *validation.QueryUser) ([]model.Category, int64, error) {
	var categories []model.Category
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	query := s.DB.WithContext(c.Context()).Model(&model.Category{}).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	result := query.Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to count categories: %+v", result.Error)
		return nil, 0, result.Error
	}

	result = query.Find(&categories)
	if result.Error != nil {
		s.Log.Errorf("Failed to get categories: %+v", result.Error)
		return nil, 0, result.Error
	}

	return categories, totalResults, nil
}

// getWeekRange returns the Monday-Sunday range for the week of t
func getWeekRange(t time.Time) (time.Time, time.Time) {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday as 7
	}
	monday := t.AddDate(0, 0, -weekday+1)
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
	sunday := monday.AddDate(0, 0, 6)
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), sunday.Location())
	return monday, sunday
}

// getMonthRange returns the first and last day of the month for t
func getMonthRange(t time.Time) (time.Time, time.Time) {
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	last := first.AddDate(0, 1, -1)
	last = time.Date(last.Year(), last.Month(), last.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), last.Location())
	return first, last
}

// getYearRange returns the first and last day of the year for t
func getYearRange(t time.Time) (time.Time, time.Time) {
	first := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	last := time.Date(t.Year(), 12, 31, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
	return first, last
}

// getDailyRange returns the start and end of the day for t
func getDailyRange(t time.Time) (time.Time, time.Time) {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
	return start, end
}

// UpsertSummary upserts daily, weekly, monthly, and yearly summaries
func UpsertSummary(db *gorm.DB, userSessionID uuid.UUID, categoryID uuid.UUID, category string, amount int64) error {
	now := time.Now()
	// Daily
	dayStart, dayEnd := getDailyRange(now)
	if err := upsertSpendingSummary(db, userSessionID, categoryID, category, amount, dayStart, dayEnd, "daily"); err != nil {
		return err
	}
	// Weekly
	weekStart, weekEnd := getWeekRange(now)
	if err := upsertSpendingSummary(db, userSessionID, categoryID, category, amount, weekStart, weekEnd, "weekly"); err != nil {
		return err
	}
	// Monthly
	monthStart, monthEnd := getMonthRange(now)
	if err := upsertSpendingSummary(db, userSessionID, categoryID, category, amount, monthStart, monthEnd, "monthly"); err != nil {
		return err
	}
	// Yearly
	yearStart, yearEnd := getYearRange(now)
	return upsertSpendingSummary(db, userSessionID, categoryID, category, amount, yearStart, yearEnd, "yearly")
}

// upsertSpendingSummary inserts or updates a summary for the given period
func upsertSpendingSummary(db *gorm.DB, userSessionID uuid.UUID, categoryID uuid.UUID, category string, amount int64, periodStart, periodEnd time.Time, periodType string) error {
	var summary model.CategorySpendingSummary
	err := db.Where("user_session_id = ? AND category_id = ? AND period_start = ? AND period_end = ? AND period_type = ?",
		userSessionID, categoryID, periodStart, periodEnd, periodType).
		First(&summary).Error

	if err == gorm.ErrRecordNotFound {
		// Insert new
		summary = model.CategorySpendingSummary{
			UserSessionID: userSessionID,
			CategoryID:    categoryID,
			Category:      category,
			TotalAmount:   amount,
			PeriodStart:   periodStart,
			PeriodEnd:     periodEnd,
			PeriodType:    periodType,
		}
		return db.Create(&summary).Error
	} else if err != nil {
		return err
	}

	// Update existing
	return db.Model(&summary).Update("total_amount", gorm.Expr("total_amount + ?", amount)).Error
}

func (s *spendingService) GetSummarySpending(c *fiber.Ctx, params *validation.QuerySpendingSummary) ([]model.CategorySpendingSummary, int64, error) {
	var summaries []model.CategorySpendingSummary
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	subQuery := s.DB.
		WithContext(c.Context()).
		Model(&model.CategorySpendingSummary{}).
		Select(
			"category_id",
			"category",
			"SUM(total_amount) AS total_amount",
			"MIN(period_start) AS period_start",
			"MAX(period_end) AS period_end",
			"period_type",
		).
		Where("period_type = ?", "daily").
		Where("user_session_id = ?", params.UserSessionID).
		Group("category_id, category, period_type")

	// Wrap the subquery into a query to count result
	countResult := s.DB.Table("(?) as summary", subQuery).Count(&totalResults)
	if countResult.Error != nil {
		s.Log.Errorf("Failed to count grouped summaries: %+v", countResult.Error)
		return nil, 0, countResult.Error
	}

	// Run the actual select
	findResult := s.DB.Raw(`
		SELECT 
			category_id, category, SUM(total_amount) AS total_amount,
			MIN(period_start) AS period_start, MAX(period_end) AS period_end,
			period_type
		FROM category_spending_summaries
		WHERE period_type = ?
		GROUP BY category_id, category, period_type
		ORDER BY total_amount DESC
	`, "daily").Scan(&summaries)

	if findResult.Error != nil {
		s.Log.Errorf("Failed to fetch grouped summaries: %+v", findResult.Error)
		return nil, 0, findResult.Error
	}

	return summaries, totalResults, nil
}
