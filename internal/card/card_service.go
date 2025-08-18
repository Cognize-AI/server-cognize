package card

import (
	"context"
	"fmt"
	"math"
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

func RebalanceCards(db *gorm.DB, listID uint, userID uint) error {
	var cards []models.Card
	if err := db.
		Where("list_id = ? AND user_id = ?", listID, userID).
		Order("card_order ASC").
		Find(&cards).Error; err != nil {
		return err
	}

	for i := range cards {
		cards[i].CardOrder = float64(i + 1)
	}

	return db.Save(&cards).Error
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

func (s *service) MoveCard(ctx context.Context, req MoveCardReq, user models.User) error {
	var prevCard, nextCard, currCard models.Card

	if req.PrevCard != 0 {
		if err := s.DB.
			Where("id = ?", req.PrevCard).
			First(&prevCard).Error; err != nil {
			return fmt.Errorf("previous card not found: %w", err)
		}
	}

	if req.NextCard != 0 {
		if err := s.DB.
			Where("id = ?", req.NextCard).
			First(&nextCard).Error; err != nil {
			return fmt.Errorf("next card not found: %w", err)
		}
	}

	if err := s.DB.
		Where("id = ?", req.CurrCard).
		First(&currCard).Error; err != nil {
		return fmt.Errorf("current card not found: %w", err)
	}
	currCard.ListID = req.ListID

	if math.Abs(nextCard.CardOrder-prevCard.CardOrder) <= 1e-9 {
		err := RebalanceCards(s.DB, req.ListID, user.ID)
		if err != nil {
			return err
		}
	}

	currCard.CardOrder = (nextCard.CardOrder + prevCard.CardOrder) / 2
	if err := s.DB.Save(&currCard).Error; err != nil {
		return err
	}

	return nil
}
