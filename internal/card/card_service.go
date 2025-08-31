package card

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/internal/tag"
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

func RebalanceCards(db *gorm.DB, listID uint) error {
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
			if err := RebalanceCards(s.DB, req.ListID); err != nil {
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
			ProfileUrl:  prospect.ProfileURL,
			AISummary:   prospect.AISummary,
		})
	}
	s.DB.Create(&cards)

	return nil, nil
}

func (s *service) GetCardByID(ctx context.Context, req GetCardByIDReq, user models.User) (*GetCardByIDResp, error) {
	var card models.Card
	var fieldVals []models.FieldValue
	var additionalContactDetails []ContactDetails
	var additionalCompanyDetails []CompanyDetails
	var cardActivity []models.Activity
	var activity []GetCardActivity
	var fieldDefIds []uint
	var fieldDefs []models.FieldDefinition

	s.DB.Preload("Tags").Preload("List").Where("id = ?", req.ID).First(&card)
	if card.ID == 0 || card.List.UserID != user.ID {
		logger.Logger.Error("card_id not found for card_id: ", zap.String("card_id", strconv.Itoa(int(req.ID))))
		return nil, errors.New("card_id not found for card_id: " + strconv.Itoa(int(req.ID)))
	}

	s.DB.Where("card_id = ?", card.ID).Find(&cardActivity)
	for _, act := range cardActivity {
		activity = append(activity, GetCardActivity{
			ID:        act.ID,
			Content:   act.Content,
			CreatedAt: act.CreatedAt,
		})
	}

	s.DB.Preload("FieldDefinition").Where("card_id = ?", card.ID).Find(&fieldVals)
	for _, fieldVal := range fieldVals {
		if models.FieldDefinitionType(fieldVal.FieldDefinition.Type) == models.CardTypeContact {
			fieldDefIds = append(fieldDefIds, fieldVal.FieldDefinition.ID)
			additionalContactDetails = append(additionalContactDetails, ContactDetails{
				ID:       fieldVal.FieldDefinition.ID,
				Name:     fieldVal.FieldDefinition.Name,
				Value:    fieldVal.Value,
				DataType: fieldVal.FieldDefinition.DataType,
			})
		} else if models.FieldDefinitionType(fieldVal.FieldDefinition.Type) == models.CardTypeCompany {
			fieldDefIds = append(fieldDefIds, fieldVal.FieldDefinition.ID)
			additionalCompanyDetails = append(additionalCompanyDetails, CompanyDetails{
				ID:       fieldVal.FieldDefinition.ID,
				Name:     fieldVal.FieldDefinition.Name,
				Value:    fieldVal.Value,
				DataType: fieldVal.FieldDefinition.DataType,
			})
		}
	}
	if len(fieldDefIds) == 0 {
		s.DB.Where("user_id = ?", user.ID).Find(&fieldDefs)
	} else {
		s.DB.Where("id not in (?) AND user_id = ?", fieldDefIds, user.ID).Find(&fieldDefs)
	}
	for _, fieldDef := range fieldDefs {
		if models.FieldDefinitionType(fieldDef.Type) == models.CardTypeContact {
			additionalContactDetails = append(additionalContactDetails, ContactDetails{
				ID:       fieldDef.ID,
				Name:     fieldDef.Name,
				Value:    "",
				DataType: fieldDef.DataType,
			})
		} else if models.FieldDefinitionType(fieldDef.Type) == models.CardTypeCompany {
			additionalCompanyDetails = append(additionalCompanyDetails, CompanyDetails{
				ID:       fieldDef.ID,
				Name:     fieldDef.Name,
				Value:    "",
				DataType: fieldDef.DataType,
			})
		}
	}

	var tags []tag.RespTag
	for _, _tag := range card.Tags {
		tags = append(tags, tag.RespTag{
			ID:    _tag.ID,
			Name:  _tag.Name,
			Color: _tag.Color,
		})
	}

	var resCard = GetCard{
		ID:          card.ID,
		Name:        card.Name,
		Designation: card.Designation,
		Email:       card.Email,
		Phone:       card.Phone,
		ImageURL:    card.ImageURL,
		ListID:      card.ListID,
		CardOrder:   card.CardOrder,
		Tags:        tags,
	}

	var res = GetCardByIDResp{
		resCard,
		card.ProfileUrl,
		card.AISummary,
		card.Location,
		card.List.Name,
		card.List.Color,
		GetCardCompanyDetails{
			Name:     card.CompanyName,
			Role:     card.CompanyRole,
			Location: card.CompanyLocation,
			Phone:    card.CompanyPhone,
			Email:    card.CompanyEmail,
		},
		additionalContactDetails,
		additionalCompanyDetails,
		activity,
	}

	return &res, nil
}

func (s *service) UpdateCardByID(ctx context.Context, req UpdateCardByIDReq, user models.User) (*UpdateCardByIDResp, error) {
	var card models.Card

	s.DB.Preload("List").Where("id = ?", req.ID).First(&card)
	if card.ID == 0 || card.ListID == 0 || card.List.UserID != user.ID {
		logger.Logger.Error("card_id not found for card_id: ", zap.String("card_id", strconv.Itoa(int(req.ID))))
		return nil, errors.New("card_id not found for card_id: " + strconv.Itoa(int(req.ID)))
	}

	card.Name = req.Name
	card.Designation = req.Designation
	card.Email = req.Email
	card.Phone = req.Phone
	card.ImageURL = req.ImageURL
	card.Location = req.Location
	card.CompanyName = req.CompanyName
	card.CompanyRole = req.CompanyRole
	card.CompanyLocation = req.CompanyLocation
	card.CompanyPhone = req.CompanyPhone
	card.CompanyEmail = req.CompanyEmail

	s.DB.Save(&card)

	return &UpdateCardByIDResp{card.ID}, nil
}
