package tag

import (
	"context"
	"errors"
	"strconv"
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

func (s *service) CreateTag(ctx context.Context, req CreateTagReq, user models.User) (*CreateTagResp, error) {
	var tag = models.Tag{
		Name:   req.Name,
		Color:  req.Color,
		UserID: user.ID,
	}

	s.DB.Create(&tag)
	return &CreateTagResp{
		tag.ID,
	}, nil
}

func (s *service) AddTag(ctx context.Context, req AddTagReq, user models.User) error {
	var card models.Card
	var tag models.Tag

	s.DB.Preload("List").Where("id = ?", req.CardID).First(&card)
	s.DB.Where("id = ?", req.TagID).First(&tag)

	if card.List.UserID != user.ID {
		logger.Logger.Error("card not exist", zap.String("card_id", strconv.Itoa(int(card.ID))), zap.String("user_id", strconv.Itoa(int(user.ID))))
		return errors.New("card not exist")
	}
	if tag.UserID != user.ID {
		logger.Logger.Error("Tag doesnt exists")
		return errors.New("tag doesnt exists")
	}

	var existingTags []models.Tag
	if err := s.DB.Model(&card).Association("Tags").Find(&existingTags); err != nil {
		return err
	}

	for _, t := range existingTags {
		if t.ID == tag.ID {
			return nil
		}
	}

	err := s.DB.Model(&card).Association("Tags").Append(&tag)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAllTags(ctx context.Context, user models.User) (*GetAllTagsResp, error) {
	var tags []models.Tag
	var respTags []RespTag

	s.DB.Where("user_id = ?", user.ID).Find(&tags)
	logger.Logger.Info("tags found: ", zap.Int("count", len(tags)))
	for _, tag := range tags {
		respTags = append(respTags, RespTag{
			tag.ID,
			tag.Name,
			tag.Color,
		})
	}

	return &GetAllTagsResp{respTags}, nil
}

func (s *service) DeleteTag(ctx context.Context, req DeleteTagReq, user models.User) error {
	var tag models.Tag
	s.DB.Where("id = ? AND user_id = ?", req.TagID, user.ID).First(&tag)
	if tag.ID == 0 {
		logger.Logger.Error("tag not exist", zap.String("tag_id", strconv.Itoa(int(tag.ID))))
		return errors.New("tag not exist")
	}

	s.DB.Delete(&tag)
	return nil
}

func (s *service) EditTag(ctx context.Context, req EditTagReq, user models.User) (*EditTagResp, error) {
	var tag models.Tag
	s.DB.Where("id = ? AND user_id = ?", req.TagID, user.ID).First(&tag)
	if tag.ID == 0 {
		logger.Logger.Error("tag not exist", zap.String("tag_id", strconv.Itoa(int(tag.ID))))
		return nil, errors.New("tag not exist")
	}

	tag.Name = req.Name
	s.DB.Save(&tag)

	return &EditTagResp{
		tag.ID,
	}, nil
}

func (s *service) RemoveTagAssociation(ctx context.Context, req RemoveTagReq, user models.User) error {
	var card models.Card
	var tag models.Tag

	g, ctx := errgroup.WithContext(ctx)

	// fetch card
	g.Go(func() error {
		return s.DB.Preload("List").Where("id = ?", req.CardID).First(&card).Error
	})

	// fetch tag
	g.Go(func() error {
		return s.DB.Where("id = ?", req.TagID).First(&tag).Error
	})

	// wait for both
	if err := g.Wait(); err != nil {
		return err
	}

	// ownership checks
	if card.List.UserID != user.ID {
		return errors.New("card not exist")
	}
	if tag.UserID != user.ID {
		return errors.New("tag doesnt exists")
	}

	if err := s.DB.Model(&card).Association("Tags").Delete(&tag); err != nil {
		return err
	}

	return nil
}
