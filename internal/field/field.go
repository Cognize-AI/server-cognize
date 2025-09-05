package field

import (
	"context"

	"github.com/Cognize-AI/client-cognize/models"
)

type CreateFieldReq struct {
	FieldName string `json:"field_name"`
	Type      string `json:"type"`
}

type CreateFieldRes struct {
	ID uint `json:"id"`
}

type InsertFieldValReq struct {
	FieldID uint   `json:"field_id"`
	CardID  uint   `json:"card_id"`
	Value   string `json:"value"`
}

type InsertFieldValRes struct {
	ID uint `json:"id"`
}

type FieldWithSample struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	SampleValue *string `json:"sample_value"`
}

type GetFieldsRes struct {
	Fields []FieldWithSample `json:"fields"`
}

type UpdateFieldDef struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	CreateField(c context.Context, req CreateFieldReq, user models.User) (*CreateFieldRes, error)
	InsertFieldVal(c context.Context, req InsertFieldValReq, user models.User) (*InsertFieldValRes, error)
	GetFields(c context.Context, user models.User) (*GetFieldsRes, error)
	UpdateFieldDefinition(c context.Context, req UpdateFieldDef, user models.User) error
}
