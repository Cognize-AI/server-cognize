package field

import (
	"context"
	"errors"
	"time"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/logger"
	"github.com/Cognize-AI/client-cognize/models"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

func (s *service) CreateField(c context.Context, req CreateFieldReq, user models.User) (*CreateFieldRes, error) {
	if !models.FieldDefinitionType(req.Type).IsFieldTypeValid() {
		logger.Logger.Error("CreateField type not valid", zap.Any("req", req))
		return nil, errors.New("CreateField type not valid")
	}

	var fieldDef models.FieldDefinition

	s.DB.Where("name = ? AND user_id = ? AND type = ?", req.FieldName, user.ID, req.Type).First(&fieldDef)
	if fieldDef.ID != 0 {
		logger.Logger.Error("Field definition already exists")
		return nil, errors.New("field definition already exists")
	}

	fieldDef = models.FieldDefinition{
		Name:   req.FieldName,
		UserID: user.ID,
		Type:   req.Type,
	}
	s.DB.Create(&fieldDef)

	return &CreateFieldRes{fieldDef.ID}, nil
}

func (s *service) InsertFieldVal(c context.Context, req InsertFieldValReq, user models.User) (*InsertFieldValRes, error) {
	var fieldDef models.FieldDefinition
	var fieldVal models.FieldValue
	var card models.Card

	g := new(errgroup.Group)

	g.Go(func() error {
		err := s.DB.
			Joins("JOIN lists ON lists.id = cards.list_id").
			Where("cards.id = ? AND lists.user_id = ?", req.CardID, user.ID).
			First(&card).Error

		if err != nil {
			logger.Logger.Error("Card not found")
			return errors.New("card not found")
		}
		return nil
	})

	g.Go(func() error {
		err := s.DB.
			Where("id = ? AND user_id = ?", req.FieldID, user.ID).
			First(&fieldDef).Error

		if err != nil || fieldDef.ID == 0 {
			logger.Logger.Error("Field definition does not exist")
			return errors.New("field definition does not exist")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	s.DB.Where("field_id = ? AND card_id = ?", req.FieldID, req.CardID).
		Assign(models.FieldValue{
			CardID:  req.CardID,
			FieldID: req.FieldID,
			Value:   req.Value,
		}).
		FirstOrCreate(&fieldVal)

	return &InsertFieldValRes{fieldVal.ID}, nil
}

func (s *service) GetFields(c context.Context, user models.User) (*GetFieldsRes, error) {
	var result []FieldWithSample

	query := `
        SELECT fd.id, fd.name, fd.type,
               (
                   SELECT fv.value
                   FROM field_values fv
                   WHERE fv.field_id = fd.id
                   LIMIT 1
               ) AS sample_value
        FROM field_definitions fd
        WHERE fd.user_id = ?
    `
	if err := s.DB.Raw(query, user.ID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &GetFieldsRes{result}, nil
}
