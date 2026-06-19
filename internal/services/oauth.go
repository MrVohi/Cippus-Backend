package services

import (
	"cippus-backend/config"
	"cippus-backend/internal/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

func NewOAuthService(cfg config.Config, db *gorm.DB) *OAuthService {
	googleConf := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	githubConf := &oauth2.Config{
		ClientID:     cfg.GithubClientID,
		ClientSecret: cfg.GithubClientSecret,
		RedirectURL:  cfg.GithubRedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	return &OAuthService{googleConf, githubConf, db, cfg.JWTSecret, cfg.FrontendURL}
}

func (s *OAuthService) GenerateState() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

func (s *OAuthService) FetchGoogleProfile(ctx context.Context, code string) (googleProfile, error) {
	token, err := s.GoogleConfig.Exchange(ctx, code)
	if err != nil {
		return googleProfile{}, err
	}
	client := s.GoogleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return googleProfile{}, err
	}
	defer resp.Body.Close()

	req := googleProfile{}
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		return googleProfile{}, err
	}
	return req, nil
}

func (s *OAuthService) FetchGithubProfile(ctx context.Context, code string) (githubProfile, error) {
	token, err := s.GithubConfig.Exchange(ctx, code)
	if err != nil {
		return githubProfile{}, err
	}
	client := s.GithubConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return githubProfile{}, err
	}
	defer resp.Body.Close()

	req := githubProfile{}
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		return githubProfile{}, err
	}

	if req.Email == "" {
		resp2, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return githubProfile{}, err
		}
		defer resp2.Body.Close()

		req2 := []githubEmail{}
		err = json.NewDecoder(resp2.Body).Decode(&req2)
		if err != nil {
			return githubProfile{}, err
		}

		for _, email := range req2 {
			if email.Primary == true && email.Verified == true {
				req.Email = email.Email
				break
			}
		}
		if req.Email == "" {
			return githubProfile{}, fmt.Errorf("no verified email on GitHub account")
		}
	}
	return req, nil
}

func (s *OAuthService) FindOrCreateUserFromOAuth(provider models.OAuthProviderName, providerUID string, email string, username string, avatarURL string) (models.User, error) {
	op := models.OAuthProvider{}
	result := s.db.Where("provider = ? AND provider_user_id = ?", provider, providerUID).First(&op)

	if result.Error == nil {
		user := models.User{}
		result = s.db.First(&user, op.UserID)
		if result.Error != nil {
			return models.User{}, result.Error
		}
		return user, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		return models.User{}, result.Error
	}

	user := models.User{}
	result = s.db.Where("email = ?", email).First(&user)

	if result.Error == nil {
		newOP := models.OAuthProvider{
			UserID:         user.ID,
			Provider:       provider,
			ProviderUserId: providerUID,
		}
		result = s.db.Create(&newOP)
		if result.Error != nil {
			return models.User{}, result.Error
		}
		return user, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		return models.User{}, result.Error
	}

	newUser := models.User{
		Email:     email,
		Username:  username,
		AvatarURL: avatarURL,
		Role:      models.RoleUser,
	}
	result = s.db.Create(&newUser)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	newOP := models.OAuthProvider{
		UserID:         newUser.ID,
		Provider:       provider,
		ProviderUserId: providerUID,
	}
	result = s.db.Create(&newOP)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return newUser, nil
}

func (s *OAuthService) IssueTokens(userID uint, role models.UserRole) (string, error) {
	accessToken, err := GenerateAccessToken(userID, role, s.secret)
	if err != nil {
		return "", err
	}

	_, tokenHash, err := GenerateRefreshToken(userID, s.secret)
	if err != nil {
		return "", err
	}

	err = StoreRefreshToken(s.db, userID, tokenHash)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
