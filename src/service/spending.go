package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SpendingService interface {
	CreateSpending(c *fiber.Ctx, req *validation.CreateSpending) (*model.Spending, error)
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

	spending := &model.Spending{
		UserSessionID: userSessionUUID, // userSessionUUID should be uuid.UUID type
		Amount:        req.Amount,
		Description:   req.Description,
		Category:      req.Category,
	}

	result := s.DB.WithContext(c.Context()).Create(spending)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Spending record already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create spending: %+v", result.Error)
	}

	return spending, result.Error
}
