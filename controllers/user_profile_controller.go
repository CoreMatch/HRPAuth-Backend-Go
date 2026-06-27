package controllers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lnb/HRPAuth-Backend-Go/database"
	"github.com/lnb/HRPAuth-Backend-Go/models"
	"github.com/lnb/HRPAuth-Backend-Go/services"
)

type UserProfileController struct{}

func NewUserProfileController() *UserProfileController {
	return &UserProfileController{}
}

type ChangeUsernameRequest struct {
	RememberToken string `json:"remember_token"`
	Username      string `json:"username"`
}

type ChangeProfileNameRequest struct {
	RememberToken string `json:"remember_token"`
	ProfileID     string `json:"profile_id"`
	Name          string `json:"name"`
}

func (uc *UserProfileController) ChangeUsername(c *gin.Context) {
	var req ChangeUsernameRequest
	token := ""
	newUsername := ""

	if err := c.ShouldBindJSON(&req); err == nil {
		token = req.RememberToken
		newUsername = req.Username
	}

	if token == "" {
		token = c.PostForm("remember_token")
	}
	if newUsername == "" {
		newUsername = c.PostForm("username")
	}

	if token == "" {
		token = c.Query("remember_token")
	}
	if newUsername == "" {
		newUsername = c.Query("username")
	}

	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "未登录或登录已过期",
		})
		return
	}

	if newUsername == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "请提供新用户名",
		})
		return
	}

	if len(newUsername) < 3 || len(newUsername) > 16 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户名长度必须在3-16个字符之间",
		})
		return
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, newUsername)
	if !matched {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户名只能包含字母、数字和下划线",
		})
		return
	}

	var user models.User
	result := database.DB.Where("remember_token = ?", token).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户不存在或token无效",
		})
		return
	}

	var count int64
	database.DB.Model(&models.User{}).
		Where("username = ? AND uid != ?", newUsername, user.UID).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "该用户名已被使用",
		})
		return
	}

	database.DB.Model(&user).Update("username", newUsername)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户名修改成功",
		"data": gin.H{
			"username": newUsername,
		},
	})
}

func (uc *UserProfileController) ChangeProfileName(c *gin.Context) {
	var req ChangeProfileNameRequest
	token := ""
	profileID := ""
	newName := ""

	if err := c.ShouldBindJSON(&req); err == nil {
		token = req.RememberToken
		profileID = req.ProfileID
		newName = req.Name
	}

	if token == "" {
		token = c.PostForm("remember_token")
	}
	if profileID == "" {
		profileID = c.PostForm("profile_id")
	}
	if newName == "" {
		newName = c.PostForm("name")
	}

	if token == "" {
		token = c.Query("remember_token")
	}
	if profileID == "" {
		profileID = c.Query("profile_id")
	}
	if newName == "" {
		newName = c.Query("name")
	}

	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "未登录或登录已过期",
		})
		return
	}

	if strings.TrimSpace(newName) == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "请提供新的角色名",
		})
		return
	}

	newName = strings.TrimSpace(newName)
	if len(newName) < 3 || len(newName) > 16 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "角色名长度必须在3-16个字符之间",
		})
		return
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, newName)
	if !matched {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "角色名只能包含字母、数字和下划线",
		})
		return
	}

	var user models.User
	result := database.DB.Where("remember_token = ?", token).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户不存在或token无效",
		})
		return
	}

	if profileID == "" {
		var profile models.Profile
		if err := database.DB.Where("user_id = ?", user.UUID).Order("created_at ASC").First(&profile).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "当前账号没有可修改的角色",
			})
			return
		}
		profileID = profile.ID
	}

	authService := services.NewAuthService()
	profile, err := authService.RenameProfile(user.UUID, profileID, newName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "角色名修改成功",
		"data": gin.H{
			"profile_id": profile.ID,
			"name":       profile.Name,
		},
	})
}