package list

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/internal/card"
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
	var resLists []CardListResponse

	s.DB.
		Where("user_id = ?", user.ID).
		Preload("Cards").
		Find(&lists)

	for _, list := range lists {
		var cards []card.GetCard

		for _, _card := range list.Cards {
			cards = append(cards, card.GetCard{
				ID:          _card.ID,
				Name:        _card.Name,
				Designation: _card.Designation,
				Email:       _card.Email,
				Phone:       _card.Phone,
				ImageURL:    _card.ImageURL,
				ListID:      _card.ListID,
				CardOrder:   _card.CardOrder,
			})
		}

		sort.Slice(cards, func(i, j int) bool {
			return cards[i].CardOrder < cards[j].CardOrder
		})

		resLists = append(resLists, CardListResponse{
			ID:        list.ID,
			Name:      list.Name,
			Color:     list.Color,
			ListOrder: list.ListOrder,
			CreatedAt: list.CreatedAt,
			UpdatedAt: list.UpdatedAt,
			Cards:     cards,
		})
	}

	sort.Slice(resLists, func(i, j int) bool {
		return resLists[i].ListOrder < resLists[j].ListOrder
	})

	return &GetListsRes{Lists: resLists}, nil
}
