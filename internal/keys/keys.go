package keys

import (
	"context"
	"time"

	"github.com/Cognize-AI/client-cognize/models"
)

type CreateAPIKeyRes struct {
	Value string `json:"value"`
}

type GetAPIKeyRes struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Service interface {
	CreateAPIKey(ctx context.Context, user models.User) (*CreateAPIKeyRes, error)
	GetAPIKey(ctx context.Context, user models.User) (*GetAPIKeyRes, error)
}
