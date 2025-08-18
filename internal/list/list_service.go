package list

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

func (s *service) CreateDefaultLists(c context.Context, user models.User) (*CreateDefaultListsRes, error) {
	var lists []models.List
	var resLists []ListResponse

	s.DB.Find(&lists).Where("user_id = ?", user.ID)
	if len(lists) > 0 {
		return nil, errors.New("default lists already exists")
	}

	lists = append(lists, models.List{
		Name:   "New Leads",
		Color:  "#F9BA0B",
		UserID: user.ID,
	})
	lists = append(lists, models.List{
		Name:   "Signed In",
		Color:  "#40C2FC",
		UserID: user.ID,
	})
	lists = append(lists, models.List{
		Name:   "Qualified",
		Color:  "#75C699",
		UserID: user.ID,
	})
	lists = append(lists, models.List{
		Name:   "Rejected",
		Color:  "#EB695B",
		UserID: user.ID,
	})

	s.DB.Create(&lists)

	for _, list := range lists {
		resLists = append(resLists, ListResponse{
			ID:        list.ID,
			Name:      list.Name,
			Color:     list.Color,
			ListOrder: list.ListOrder,
			CreatedAt: list.CreatedAt,
			UpdatedAt: list.UpdatedAt,
		})
	}
	logger.Logger.Info("Created default lists")
	return &CreateDefaultListsRes{Lists: resLists}, nil
}

func (s *service) GetLists(c context.Context, user models.User) (*GetListsRes, error) {
	var lists []models.List
	var resLists []ListResponse

	s.DB.Find(&lists).Where("user_id = ?", user.ID)

	for _, list := range lists {
		resLists = append(resLists, ListResponse{
			ID:        list.ID,
			Name:      list.Name,
			Color:     list.Color,
			ListOrder: list.ListOrder,
			CreatedAt: list.CreatedAt,
			UpdatedAt: list.UpdatedAt,
		})
	}

	return &GetListsRes{Lists: resLists}, nil
}
