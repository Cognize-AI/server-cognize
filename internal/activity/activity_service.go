package activity

import (
	"context"
	"errors"
	"time"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/logger"
	"github.com/Cognize-AI/client-cognize/models"
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

func (s *service) CreateActivity(ctx context.Context, req CreateActivityReq, user models.User) (*CreateActivityResp, error) {
	var activity models.Activity
	var card models.Card

	s.DB.Preload("List").Where("id = ?", req.CardID).First(&card)
	if card.ID == 0 || card.List.UserID != user.ID {
		logger.Logger.Error("Card not found")
		return nil, errors.New("card not found")
	}

	activity = models.Activity{
		Content: req.Text,
		CardID:  req.CardID,
	}
	s.DB.Create(&activity)

	return &CreateActivityResp{activity.ID}, nil
}

func (s *service) DeleteActivity(ctx context.Context, req DeleteActivityReq, user models.User) error {
	var activity models.Activity

	if err := s.DB.Preload("Card.List").Where("id = ?", req.ID).First(&activity).Error; err != nil {
		logger.Logger.Error("Activity not found")
		return errors.New("activity not found")
	}

	if activity.Card.ID == 0 || activity.Card.List.UserID != user.ID {
		logger.Logger.Error("Unauthorized or card not found")
		return errors.New("card not found")
	}

	if err := s.DB.Delete(&activity).Error; err != nil {
		logger.Logger.Error("Failed to delete activity", zap.Error(err))
		return errors.New("failed to delete activity")
	}

	return nil
}

func (s *service) UpdateActivity(ctx context.Context, req UpdateActivityReq, user models.User) (*UpdateActivityResp, error) {
	var activity models.Activity

	if err := s.DB.Preload("Card.List").Where("id = ?", req.ID).First(&activity).Error; err != nil {
		logger.Logger.Error("Activity not found")
		return nil, errors.New("activity not found")
	}

	if activity.Card.ID == 0 || activity.Card.List.UserID != user.ID {
		logger.Logger.Error("Unauthorized or card not found")
		return nil, errors.New("card not found")
	}

	activity.Content = req.Text
	if err := s.DB.Save(&activity).Error; err != nil {
		logger.Logger.Error("Failed to update activity", zap.Error(err))
		return nil, errors.New("failed to update activity")
	}

	return &UpdateActivityResp{activity.ID}, nil
}
