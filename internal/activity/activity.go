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

type Service interface {
	CreateActivity(ctx context.Context, req CreateActivityReq, user models.User) (*CreateActivityResp, error)
}
