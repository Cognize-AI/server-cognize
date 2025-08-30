package activity

import (
	"context"

	"github.com/Cognize-AI/client-cognize/models"
)

type CreateActivityReq struct {
	CardID uint   `json:"card_id"`
	Text   string `json:"text"`
}

type CreateActivityResp struct {
	ID uint `json:"id"`
}

type DeleteActivityReq struct {
	ID uint `uri:"id" binding:"required"`
}

type UpdateActivityReq struct {
	ID   uint   `uri:"id" binding:"required"`
	Text string `json:"text"`
}

type UpdateActivityResp struct {
	ID uint `json:"id"`
}

type Service interface {
	CreateActivity(ctx context.Context, req CreateActivityReq, user models.User) (*CreateActivityResp, error)
	DeleteActivity(ctx context.Context, req DeleteActivityReq, user models.User) error
	UpdateActivity(ctx context.Context, req UpdateActivityReq, user models.User) (*UpdateActivityResp, error)
}
