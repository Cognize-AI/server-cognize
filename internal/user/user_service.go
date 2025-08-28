package user

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Cognize-AI/server-cognize/config"
	"github.com/Cognize-AI/server-cognize/logger"
	"github.com/Cognize-AI/server-cognize/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	timeout time.Duration
	DB      *gorm.DB
}

func NewService() Service {
	return &service{
		time.Duration(20) * time.Second,
		config.DB,
	}
}

func (s *service) Me(c context.Context, user models.User) (*GetMeRes, error) {
	s.DB.First(&user, "id=?", user.ID)
	if user.ID == 0 {
		logger.Logger.Error("user not found: ", zap.String("id", strconv.Itoa(int(user.ID))))
		return nil, errors.New("user not found")
	}

	res := &GetMeRes{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
	}

	return res, nil
}
