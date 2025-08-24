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
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type GetAllTagsResp struct {
	Tags []RespTag `json:"tags"`
}

type DeleteTagReq struct {
	TagID uint `uri:"id"`
}

type EditTagReq struct {
	TagID uint   `json:"id"`
	Name  string `json:"name"`
}

type EditTagResp struct {
	ID uint `json:"id"`
}

//🌸 Pastel Pink → #F8BBD0
//
//🌿 Mint Green → #B2EBF2
//
//🌼 Soft Yellow → #FFF9C4
//
//🌊 Baby Blue → #BBDEFB
//
//🍑 Peach → #FFE0B2

type Service interface {
	CreateTag(ctx context.Context, req CreateTagReq, user models.User) (*CreateTagResp, error)
	AddTag(ctx context.Context, req AddTagReq, user models.User) error
	GetAllTags(ctx context.Context, user models.User) (*GetAllTagsResp, error)
	DeleteTag(ctx context.Context, req DeleteTagReq, user models.User) error
	EditTag(ctx context.Context, req EditTagReq, user models.User) (*EditTagResp, error)
}
