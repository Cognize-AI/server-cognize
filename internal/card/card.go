package card

import "context"

type CreateCardReq struct {
	Name        string
	Designation string
	Email       string
	Phone       string
	ImageURL    string
	ListID      uint `gorm:"index"`
}

type CreateCardResp struct {
	ID uint `json:"id"`
}

type Service interface {
	CreateCard(ctx context.Context, req CreateCardReq) (*CreateCardResp, error)
}
