package activity

import (
	"context"
	"errors"
	"time"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/logger"
	"github.com/Cognize-AI/client-cognize/models"
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
