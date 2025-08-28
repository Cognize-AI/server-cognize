package list

import (
	"context"
	"time"

	"github.com/Cognize-AI/server-cognize/internal/card"
	"github.com/Cognize-AI/server-cognize/models"
)

type GetListResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	ListOrder float64   `json:"list_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CardListResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Color     string         `json:"color"`
	ListOrder float64        `json:"list_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Cards     []card.GetCard `json:"cards"`
}

type CreateDefaultListsRes struct {
	Lists []GetListResponse `json:"lists"`
}

type GetListsRes struct {
	Lists []CardListResponse `json:"lists"`
}

type Service interface {
	CreateDefaultLists(c context.Context, user models.User) (*CreateDefaultListsRes, error)
	GetLists(c context.Context, user models.User) (*GetListsRes, error)
}
