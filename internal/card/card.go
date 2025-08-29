package card

import (
	"context"
	"time"

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
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	DataType string `json:"data_type"`
}

type CompanyDetails struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	DataType string `json:"data_type"`
}

type GetCardCompanyDetails struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Location string `json:"location"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type GetCardActivity struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type GetCardByIDResp struct {
	GetCard
	Location          string                `json:"location"`
	ListName          string                `json:"list_name"`
	ListColor         string                `json:"list_color"`
	Company           GetCardCompanyDetails `json:"company"`
	AdditionalContact []ContactDetails      `json:"additional_contact"`
	AdditionalCompany []CompanyDetails      `json:"additional_company"`
	Activity          []GetCardActivity     `json:"activity"`
}

type BulkCreateResp struct {
}

type UpdateCardByIDReq struct {
	ID              uint   `uri:"id" binding:"required"`
	Name            string `json:"name"`
	Designation     string `json:"designation"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	ImageURL        string `json:"image_url"`
	Location        string `json:"location"`
	CompanyName     string `json:"company_name"`
	CompanyRole     string `json:"company_role"`
	CompanyLocation string `json:"company_location"`
	CompanyPhone    string `json:"company_phone"`
	CompanyEmail    string `json:"company_email"`
}

type UpdateCardByIDResp struct {
	ID uint `json:"id"`
}

type Service interface {
	CreateCard(ctx context.Context, req CreateCardReq, user models.User) (*CreateCardResp, error)
	MoveCard(ctx context.Context, req MoveCardReq, user models.User) error
	DeleteCard(ctx context.Context, req DeleteCardReq, user models.User) (*DeleteCardResp, error)
	UpdateCard(ctx context.Context, req UpdateCardReq, user models.User) (*UpdateCardResp, error)
	BulkCreate(ctx context.Context, req BulkCreateReq, key models.Key) (*BulkCreateResp, error)
	GetCardByID(ctx context.Context, req GetCardByIDReq, user models.User) (*GetCardByIDResp, error)
	UpdateCardByID(ctx context.Context, req UpdateCardByIDReq, user models.User) (*UpdateCardByIDResp, error)
}
