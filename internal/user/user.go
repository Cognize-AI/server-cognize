package user

import (
	"context"

	"github.com/Cognize-AI/client-cognize/models"
)

type GetMeRes struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}

type Service interface {
	Me(c context.Context, user models.User) (*GetMeRes, error)
}
