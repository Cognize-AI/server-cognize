package keys

import (
	"context"
	"time"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/logger"
	"github.com/Cognize-AI/client-cognize/models"
	"github.com/Cognize-AI/client-cognize/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	timeout time.Duration
	DB      *gorm.DB
}

func NewService() Service {
	return &service{
		time.Duration(20) * time.Second,
		config.DB,
	}
}

func (s *service) CreateAPIKey(ctx context.Context, user models.User) (*CreateAPIKeyRes, error) {
	var key models.Key

	err := s.DB.Where("user_id = ? AND name = ?", user.ID, "API").First(&key).Error
	if err == nil && key.ID != 0 {
		logger.Logger.Warn("api key already exists")
		return &CreateAPIKeyRes{key.Value}, nil
	}

	logger.Logger.Info("creating api key")
	value, err := util.GenerateAPIKey()
	if err != nil {
		logger.Logger.Error("failed to generate api key", zap.Error(err))
		return nil, err
	}

	key = models.Key{
		Name:   "API",
		Value:  value,
		UserID: user.ID,
	}

	if err := s.DB.Create(&key).Error; err != nil {
		logger.Logger.Error("failed to save api key", zap.Error(err))
		return nil, err
	}

	return &CreateAPIKeyRes{key.Value}, nil
}

func (s *service) GetAPIKey(ctx context.Context, user models.User) (*GetAPIKeyRes, error) {
	var key models.Key

	err := s.DB.Where("user_id = ? AND name = ?", user.ID, "API").First(&key).Error
	if err != nil {
		logger.Logger.Error("failed to get api key", zap.Error(err))
		return nil, err
	}

	return &GetAPIKeyRes{
		ID:        key.ID,
		Key:       key.Value,
		Name:      key.Name,
		CreatedAt: key.CreatedAt,
	}, nil
}
