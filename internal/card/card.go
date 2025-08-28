package card

import (
	"context"

	"github.com/Cognize-AI/client-cognize/internal/tag"
	"github.com/Cognize-AI/client-cognize/models"
)

type GetCard struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Designation string        `json:"designation"`
	Email       string        `json:"email"`
	Phone       string        `json:"phone"`
	ImageURL    string        `json:"image_url"`
	ListID      uint          `gorm:"index"`
	CardOrder   float64       `gorm:"autoIncrement"`
	Tags        []tag.RespTag `json:"tags"`
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

type DeleteCardReq struct {
	ID uint `uri:"id" binding:"required"`
}

type DeleteCardResp struct {
	ID uint `json:"id"`
}

type UpdateCardReq struct {
	ID          uint   `uri:"id" binding:"required"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ImageURL    string `json:"image_url"`
	ListID      uint   `json:"list_id"`
}

type UpdateCardResp struct {
	ID uint `uri:"id"`
}

type BulkProspect struct {
	Name        string `json:"name"`
	Designation string `json:"designation"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ImageURL    string `json:"image_url"`
}

type BulkCreateReq struct {
	ListID    uint           `json:"list_id"`
	Prospects []BulkProspect `json:"prospects"`
}

type GetCardByIDReq struct {
	ID uint `uri:"id" binding:"required"`
}

type ContactDetails struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	DataType string `json:"data_type"`
}

type CompanyDetails struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	DataType string `json:"data_type"`
}

type GetCardByIDResp struct {
	GetCard
	ListName          string           `json:"list_name"`
	AdditionalContact []ContactDetails `json:"additional_contact"`
	AdditionalCompany []CompanyDetails `json:"additional_company"`
}

type BulkCreateResp struct {
}

type Service interface {
	CreateCard(ctx context.Context, req CreateCardReq, user models.User) (*CreateCardResp, error)
	MoveCard(ctx context.Context, req MoveCardReq, user models.User) error
	DeleteCard(ctx context.Context, req DeleteCardReq, user models.User) (*DeleteCardResp, error)
	UpdateCard(ctx context.Context, req UpdateCardReq, user models.User) (*UpdateCardResp, error)
	BulkCreate(ctx context.Context, req BulkCreateReq, key models.Key) (*BulkCreateResp, error)
	GetCardByID(ctx context.Context, req GetCardByIDReq, user models.User) (*GetCardByIDResp, error)
}
