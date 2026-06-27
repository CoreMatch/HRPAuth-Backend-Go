package controllers

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lnb/HRPAuth-Backend-Go/database"
	"github.com/lnb/HRPAuth-Backend-Go/models"
	"github.com/lnb/HRPAuth-Backend-Go/services"
)

type TextureController struct {
	textureService *services.TextureService
	authService    *services.AuthService
}

func NewTextureController() *TextureController {
	return &TextureController{
		textureService: services.NewTextureService(),
		authService:    services.NewAuthService(),
	}
}

type UploadTextureRequest struct {
	RememberToken string `json:"remember_token"`
	ProfileID     string `json:"profile_id"`
	TextureType   string `json:"texture_type"`
	Model         string `json:"model"`
}

func (tc *TextureController) UploadTexture(c *gin.Context) {
	token := ""
	profileID := ""
	textureType := ""
	model := ""

	contentType := c.ContentType()
	if strings.Contains(contentType, "application/json") {
		var req UploadTextureRequest
		if err := c.ShouldBindJSON(&req); err == nil {
			token = req.RememberToken
			profileID = req.ProfileID
			textureType = req.TextureType
			model = req.Model
		}
	}

	if token == "" {
		token = c.PostForm("remember_token")
	}
	if profileID == "" {
		profileID = c.PostForm("profile_id")
	}
	if textureType == "" {
		textureType = c.PostForm("texture_type")
	}
	if model == "" {
		model = c.PostForm("model")
	}

	if token == "" {
		token = c.Query("remember_token")
	}
	if profileID == "" {
		profileID = c.Query("profile_id")
	}
	if textureType == "" {
		textureType = c.Query("texture_type")
	}
	if model == "" {
		model = c.Query("model")
	}

	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "未登录或登录已过期",
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

	if textureType != "skin" && textureType != "cape" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的材质类型，只能是 skin 或 cape",
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

	if !tc.authService.IsProfileOwnedByUser(profileID, user.UUID) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无权操作该角色",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "未上传文件",
		})
		return
	}

	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无法读取文件",
		})
		return
	}
	defer fileData.Close()

	data, err := io.ReadAll(fileData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文件读取失败",
		})
		return
	}

	if err := tc.textureService.UploadTextureByUser(user.UUID, profileID, textureType, model, data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "材质上传成功",
		"data": gin.H{
			"profile_id":   profileID,
			"texture_type": textureType,
		},
	})
}

type DeleteTextureRequest struct {
	RememberToken string `json:"remember_token"`
	ProfileID     string `json:"profile_id"`
	TextureType   string `json:"texture_type"`
}

func (tc *TextureController) DeleteTexture(c *gin.Context) {
	token := ""
	profileID := ""
	textureType := ""

	contentType := c.ContentType()
	if strings.Contains(contentType, "application/json") {
		var req DeleteTextureRequest
		if err := c.ShouldBindJSON(&req); err == nil {
			token = req.RememberToken
			profileID = req.ProfileID
			textureType = req.TextureType
		}
	}

	if token == "" {
		token = c.PostForm("remember_token")
	}
	if profileID == "" {
		profileID = c.PostForm("profile_id")
	}
	if textureType == "" {
		textureType = c.PostForm("texture_type")
	}

	if token == "" {
		token = c.Query("remember_token")
	}
	if profileID == "" {
		profileID = c.Query("profile_id")
	}
	if textureType == "" {
		textureType = c.Query("texture_type")
	}

	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "未登录或登录已过期",
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

	if textureType != "skin" && textureType != "cape" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的材质类型，只能是 skin 或 cape",
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

	if !tc.authService.IsProfileOwnedByUser(profileID, user.UUID) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无权操作该角色",
		})
		return
	}

	if err := tc.textureService.RemoveTextureByUser(user.UUID, profileID, textureType); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "材质删除成功",
		"data": gin.H{
			"profile_id":   profileID,
			"texture_type": textureType,
		},
	})
}

type GetTextureRequest struct {
	RememberToken string `json:"remember_token"`
	ProfileID     string `json:"profile_id"`
}

type TextureResponse struct {
	TextureType string `json:"texture_type"`
	URL         string `json:"url"`
	Model       string `json:"model,omitempty"`
}

func (tc *TextureController) GetTexture(c *gin.Context) {
	token := ""
	profileID := ""

	contentType := c.ContentType()
	if strings.Contains(contentType, "application/json") {
		var req GetTextureRequest
		if err := c.ShouldBindJSON(&req); err == nil {
			token = req.RememberToken
			profileID = req.ProfileID
		}
	}

	if token == "" {
		token = c.PostForm("remember_token")
	}
	if profileID == "" {
		profileID = c.PostForm("profile_id")
	}

	if token == "" {
		token = c.Query("remember_token")
	}
	if profileID == "" {
		profileID = c.Query("profile_id")
	}

	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "未登录或登录已过期",
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
				"message": "当前账号没有角色",
			})
			return
		}
		profileID = profile.ID
	}

	if !tc.authService.IsProfileOwnedByUser(profileID, user.UUID) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无权查看该角色",
		})
		return
	}

	skinInfo, _ := tc.textureService.GetTextureByProfile(profileID, "skin")
	capeInfo, _ := tc.textureService.GetTextureByProfile(profileID, "cape")

	textures := make([]TextureResponse, 0)
	if skinInfo != nil {
		skinResp := TextureResponse{
			TextureType: "skin",
			URL:         skinInfo.URL,
		}
		if skinInfo.Metadata != nil {
			if model, ok := skinInfo.Metadata["model"]; ok {
				skinResp.Model = model.(string)
			}
		}
		textures = append(textures, skinResp)
	}
	if capeInfo != nil {
		textures = append(textures, TextureResponse{
			TextureType: "cape",
			URL:         capeInfo.URL,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取材质信息成功",
		"data": gin.H{
			"profile_id": profileID,
			"textures":   textures,
		},
	})
}