package list

import (
	"context"
	"time"

	"github.com/Cognize-AI/client-cognize/models"
)

type ListResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	ListOrder uint      `json:"list_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDefaultListsRes struct {
	Lists []ListResponse `json:"lists"`
}

type GetListsRes struct {
	Lists []ListResponse `json:"lists"`
}

type Service interface {
	CreateDefaultLists(c context.Context, user models.User) (*CreateDefaultListsRes, error)
	GetLists(c context.Context, user models.User) (*GetListsRes, error)
}
