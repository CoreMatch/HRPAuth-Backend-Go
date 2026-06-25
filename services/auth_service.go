package services

import (
	"time"

	"github.com/lnb/HRPAuth-Backend-Go/config"
	"github.com/lnb/HRPAuth-Backend-Go/database"
	"github.com/lnb/HRPAuth-Backend-Go/models"
	"github.com/lnb/HRPAuth-Backend-Go/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

type UserInfo struct {
	UUID     string
	Email    string
	Username string
	Password string
	Locale   string
}

func (as *AuthService) VerifyCredentials(identifier, password string, allowUsernameLogin bool) *UserInfo {
	nonEmailLogin := allowUsernameLogin
	if !allowUsernameLogin {
		nonEmailLogin = config.AppConfig.Yggdrasil.FeatureFlags.NonEmailLogin
	}

	var user models.User
	var err error

	if nonEmailLogin {
		err = database.DB.
			Joins("JOIN profiles ON users.uuid = profiles.user_id").
			Where("users.email = ? OR profiles.name = ?", identifier, identifier).
			First(&user).Error
	} else {
		err = database.DB.Where("email = ?", identifier).First(&user).Error
	}

	if err != nil {
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil
	}

	return &UserInfo{
		UUID:     user.UUID,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		Locale:   user.Locale,
	}
}

type ProfileInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Model string `json:"model,omitempty"`
}

func (as *AuthService) GetUserProfiles(userUUID string) []ProfileInfo {
	var profiles []models.Profile
	database.DB.Where("user_id = ?", userUUID).Find(&profiles)

	result := make([]ProfileInfo, 0, len(profiles))
	for _, p := range profiles {
		pi := ProfileInfo{
			ID:   p.ID,
			Name: p.Name,
		}
		if p.Model != "" {
			pi.Model = p.Model
		}
		result = append(result, pi)
	}
	return result
}

func (as *AuthService) CreateToken(accessToken, clientToken, userID, profileID string, expiresInDays int) bool {
	token := models.Token{
		AccessToken:       accessToken,
		ClientToken:       clientToken,
		UserID:            userID,
		SelectedProfileID: profileID,
		IssuedAt:          utils.CurrentTimestampMillis(),
		ExpiresInDays:     expiresInDays,
		State:             "valid",
	}
	result := database.DB.Create(&token)
	return result.Error == nil
}

func (as *AuthService) InvalidateToken(accessToken string) bool {
	result := database.DB.Model(&models.Token{}).
		Where("access_token = ?", accessToken).
		Update("state", "invalid")
	return result.Error == nil
}

func (as *AuthService) ValidateToken(accessToken string, clientToken string) *models.Token {
	var token models.Token
	result := database.DB.Where("access_token = ? AND state = ?", accessToken, "valid").First(&token)
	if result.Error != nil {
		return nil
	}

	if clientToken != "" && clientToken != token.ClientToken {
		return nil
	}

	return &token
}

func (as *AuthService) GetProfileByID(profileID string) *ProfileInfo {
	var profile models.Profile
	var user models.User

	result := database.DB.Where("id = ?", profileID).First(&profile)
	if result.Error != nil {
		return nil
	}

	database.DB.Where("uuid = ?", profile.UserID).First(&user)

	return &ProfileInfo{
		ID:    profile.ID,
		Name:  user.Username,
		Model: profile.Model,
	}
}

func (as *AuthService) IsProfileOwnedByUser(profileID, userID string) bool {
	var count int64
	database.DB.Model(&models.Profile{}).
		Where("id = ? AND user_id = ?", profileID, userID).
		Count(&count)
	return count > 0
}

func (as *AuthService) CreateSession(profileID, serverID, ip string) bool {
	session := models.Session{
		ProfileID: profileID,
		ServerID:  serverID,
		IP:        ip,
		ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.Yggdrasil.Security.SessionExpirySeconds) * time.Second),
	}
	result := database.DB.Create(&session)
	return result.Error == nil
}

func (as *AuthService) GetSessionByProfileAndServer(profileName, serverID string) *models.Session {
	var profile models.Profile
	database.DB.Where("name = ?", profileName).First(&profile)
	if profile.ID == "" {
		return nil
	}

	var session models.Session
	result := database.DB.
		Where("profile_id = ? AND server_id = ? AND expires_at > ?", profile.ID, serverID, time.Now()).
		Order("created_at DESC").
		First(&session)

	if result.Error != nil {
		return nil
	}
	return &session
}

func (as *AuthService) GetProfileByName(name string) *models.Profile {
	var profile models.Profile
	result := database.DB.Where("name = ?", name).First(&profile)
	if result.Error != nil {
		return nil
	}
	return &profile
}
