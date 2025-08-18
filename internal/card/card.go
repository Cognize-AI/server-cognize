package card

import (
	"context"

	"github.com/Cognize-AI/client-cognize/models"
)

type GetCard struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Designation string  `json:"designation"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	ImageURL    string  `json:"image_url"`
	ListID      uint    `gorm:"index"`
	CardOrder   float64 `gorm:"autoIncrement"`
}

type CreateCardReq struct {
	Name        string `json:"name"`
	Designation string `json:"designation"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ImageURL    string `json:"image_url"`
	ListID      uint   `json:"list_id"`
}

type CreateCardResp struct {
	ID uint `json:"id"`
}

type MoveCardReq struct {
	PrevCard uint `json:"prev_card"`
	CurrCard uint `json:"curr_card"`
	NextCard uint `json:"next_card"`
	ListID   uint `json:"list_id"`
}

type Service interface {
	CreateCard(ctx context.Context, req CreateCardReq, user models.User) (*CreateCardResp, error)
	MoveCard(ctx context.Context, req MoveCardReq, user models.User) error
}
