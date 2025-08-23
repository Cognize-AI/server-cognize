package tag

import (
	"context"

	"github.com/Cognize-AI/client-cognize/models"
)

type CreateTagReq struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CreateTagResp struct {
	ID uint `json:"id"`
}

type AddTagReq struct {
	TagID  uint `json:"tag_id"`
	CardID uint `json:"card_id"`
}

type RespTag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type GetAllTagsResp struct {
	Tags []RespTag `json:"tags"`
}

type Service interface {
	CreateTag(ctx context.Context, req CreateTagReq, user models.User) (*CreateTagResp, error)
	AddTag(ctx context.Context, req AddTagReq, user models.User) error
	GetAllTags(ctx context.Context, user models.User) (*GetAllTagsResp, error)
}
