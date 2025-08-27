package card

import (
	"context"
	"errors"
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
		Where("list_id = ?", listID).
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

	// Load prev card if provided
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

	// Decide new order
	if req.PrevCard != 0 && req.NextCard != 0 {
		// Case 1: Move between two cards
		if math.Abs(nextCard.CardOrder-prevCard.CardOrder) <= 1e-9 {
			// Rebalance list if gap is too small
			if err := RebalanceCards(s.DB, req.ListID, user.ID); err != nil {
				return fmt.Errorf("failed to rebalance cards: %w", err)
			}
			// Reload neighbors after rebalance
			if err := s.DB.Where("id = ?", req.PrevCard).First(&prevCard).Error; err != nil {
				return fmt.Errorf("previous card not found after rebalance: %w", err)
			}
			if err := s.DB.Where("id = ?", req.NextCard).First(&nextCard).Error; err != nil {
				return fmt.Errorf("next card not found after rebalance: %w", err)
			}
		}
		currCard.CardOrder = (nextCard.CardOrder + prevCard.CardOrder) / 2

	} else if req.PrevCard == 0 && req.NextCard != 0 {
		// Case 2: Insert at beginning
		currCard.CardOrder = nextCard.CardOrder - 1

	} else if req.NextCard == 0 && req.PrevCard != 0 {
		// Case 3: Insert at end
		currCard.CardOrder = prevCard.CardOrder + 1

	} else {
		// Case 4: Only card in the list
		currCard.CardOrder = 1
	}

	if err := s.DB.Save(&currCard).Error; err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}

	return nil
}

func (s *service) DeleteCard(ctx context.Context, req DeleteCardReq, user models.User) (*DeleteCardResp, error) {
	var card models.Card
	s.DB.Preload("List").Where("id = ?", req.ID).First(&card)
	if card.ID == 0 {
		logger.Logger.Error("list not found for card_id: ", zap.String("card_id", strconv.Itoa(int(req.ID))))
		return nil, errors.New("list not found for card_id: " + strconv.Itoa(int(req.ID)))
	}
	if card.List.UserID != user.ID {
		logger.Logger.Error("card_id not found for user_id: ", zap.String("user_id", strconv.Itoa(int(user.ID))))
		return nil, errors.New("card not found for user")
	}

	if err := s.DB.Delete(&card).Error; err != nil {
		logger.Logger.Error("Error deleting card", zap.Error(err))
		return nil, fmt.Errorf("failed to delete card: %w", err)
	}
	return &DeleteCardResp{
		card.ID,
	}, nil
}

func (s *service) UpdateCard(ctx context.Context, req UpdateCardReq, user models.User) (*UpdateCardResp, error) {
	var card models.Card
	s.DB.Preload("List").Where("id = ?", req.ID).First(&card)
	if card.ID == 0 {
		logger.Logger.Error("card not found for card_id: ", zap.String("card_id", strconv.Itoa(int(req.ID))))
		return nil, errors.New("card not found for card_id: " + strconv.Itoa(int(req.ID)))
	}
	if card.List.ID == 0 {
		logger.Logger.Error("list not found for card_id: ", zap.String("card_id", strconv.Itoa(int(req.ID))))
		return nil, errors.New("list not found for card_id: " + strconv.Itoa(int(req.ID)))
	}
	if card.List.UserID != user.ID {
		logger.Logger.Error("card_id not found for user_id: ", zap.String("user_id", strconv.Itoa(int(user.ID))))
		return nil, errors.New("card not found for user")
	}

	card.Name = req.Name
	card.Designation = req.Designation
	card.Email = req.Email
	card.Phone = req.Phone
	card.ImageURL = req.ImageURL

	s.DB.Save(&card)
	return &UpdateCardResp{card.ID}, nil
}

func (s *service) BulkCreate(ctx context.Context, req BulkCreateReq, key models.Key) (*BulkCreateResp, error) {
	var cards []models.Card
	var list models.List

	s.DB.Where("id = ? AND user_id = ?", req.ListID, key.UserID).First(&list)
	if list.ID == 0 {
		logger.Logger.Error("list not found for list_id: ", zap.String("list_id", strconv.Itoa(int(req.ListID))))
		return nil, errors.New("list not found for list_id: " + strconv.Itoa(int(req.ListID)))
	}

	var maxOrder float64
	s.DB.Model(&models.Card{}).Select("COALESCE(MAX(card_order), 0)").Scan(&maxOrder)

	for i, prospect := range req.Prospects {
		cards = append(cards, models.Card{
			Name:        prospect.Name,
			Designation: prospect.Designation,
			Email:       prospect.Email,
			Phone:       prospect.Phone,
			ImageURL:    prospect.ImageURL,
			ListID:      req.ListID,
			CardOrder:   maxOrder + float64(i+1),
		})
	}
	s.DB.Create(&cards)
	
	return nil, nil
}
