package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
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

func (s *service) GetRedirectURL(c context.Context) (*GetRedirectURLResp, error) {
	url := config.GoogleOAuthConfig.AuthCodeURL("random-state-token", oauth2.AccessTypeOffline)
	return &GetRedirectURLResp{
		RedirectURL: url,
	}, nil
}

func (s *service) HandleGoogleCallback(c context.Context, req *HandleGoogleCallbackReq) (*HandleGoogleCallbackResp, error) {
	_config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), req.Code)
	if err != nil {
		return nil, errors.New("code exchange failed")
	}

	client := config.GoogleOAuthConfig.Client(context.Background(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.New("failed to get user info")
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var googleUser GoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		log.Fatal(err)
	}

	var user models.User
	s.DB.Where("email = ?", googleUser.Email).First(&user)
	if user.ID == 0 {
		user = models.User{
			Name:           googleUser.Name,
			Email:          googleUser.Email,
			Password:       "",
			ProfilePicture: googleUser.Picture,
		}
		s.DB.Create(&user)
		var lists []models.List

		lists = append(lists, models.List{
			Name:   "New Leads",
			Color:  "#F9BA0B",
			UserID: user.ID,
		})
		lists = append(lists, models.List{
			Name:   "Follow Up",
			Color:  "#40C2FC",
			UserID: user.ID,
		})
		lists = append(lists, models.List{
			Name:   "Qualified",
			Color:  "#75C699",
			UserID: user.ID,
		})
		lists = append(lists, models.List{
			Name:   "Rejected",
			Color:  "#EB695B",
			UserID: user.ID,
		})

		s.DB.Create(&lists)

		var tags []models.Tag
		tags = append(tags, models.Tag{
			Name:   "ux researcher",
			Color:  "#A78BFA",
			UserID: user.ID,
		}, models.Tag{
			Name:   "product designer",
			Color:  "#FCA5A5",
			UserID: user.ID,
		}, models.Tag{
			Name:   "content strategist",
			Color:  "#34D399",
			UserID: user.ID,
		}, models.Tag{
			Name:   "SEO specialist",
			Color:  "#60A5FA",
			UserID: user.ID,
		}, models.Tag{
			Name:   "brand strategist",
			Color:  "#FBBF24",
			UserID: user.ID,
		})
		s.DB.Create(&tags)
	} else {
		user.ProfilePicture = googleUser.Picture
		s.DB.Save(&user)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := jwtToken.SignedString([]byte(_config.JwtSecret))
	if err != nil {
		return nil, err
	}

	return &HandleGoogleCallbackResp{
		Token:          tokenString,
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
	}, nil
}
