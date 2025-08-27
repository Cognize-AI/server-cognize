package keys

import (
	"context"

	"github.com/Cognize-AI/client-cognize/models"
)

type CreateAPIKeyRes struct {
	Value string `json:"value"`
}

type Service interface {
	CreateAPIKey(ctx context.Context, user models.User) (*CreateAPIKeyRes, error)
}
