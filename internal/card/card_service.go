package card

import (
	"context"
	"strconv"
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

func (s *service) CreateCard(ctx context.Context, req CreateCardReq, user models.User) (*CreateCardResp, error) {
	var list models.List
	s.DB.Where("id = ? AND user_id = ?", req.ListID, user.ID).First(&list)
	if list.ID == 0 {
		logger.Logger.Error("list not found", zap.String("list_id", strconv.Itoa(int(req.ListID))))
	}

	var maxOrder float64
	s.DB.Model(&models.Card{}).Select("COALESCE(MAX(card_order), 0)").Scan(&maxOrder)

	var card = models.Card{
		Name:        req.Name,
		Designation: req.Designation,
		Email:       req.Email,
		Phone:       req.Phone,
		ImageURL:    req.ImageURL,
		ListID:      req.ListID,
		CardOrder:   maxOrder + 1,
	}
	s.DB.Create(&card)

	return &CreateCardResp{card.ID}, nil
}
