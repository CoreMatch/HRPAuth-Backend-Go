package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lnb/HRPAuth-Backend-Go/config"
	"github.com/lnb/HRPAuth-Backend-Go/services"
)

type CaptchaController struct {
	captchaService *services.CaptchaService
}

func NewCaptchaController() *CaptchaController {
	return &CaptchaController{
		captchaService: services.NewCaptchaService(),
	}
}

// Status reports whether captcha verification is enabled. Returns 1 when
// enabled, 0 when disabled.
func (cc *CaptchaController) Status(c *gin.Context) {
	enabled := 0
	if config.AppConfig.Security.EnableCaptcha {
		enabled = 1
	}
	c.JSON(http.StatusOK, gin.H{
		"enabled": enabled,
	})
}

// Generate issues a new captcha. Returns the token and the URL at which the
// PNG image can be fetched. The code itself is never returned to clients.
func (cc *CaptchaController) Generate(c *gin.Context) {
	if !config.AppConfig.Security.EnableCaptcha {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Captcha is disabled",
		})
		return
	}

	token, _, err := cc.captchaService.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate captcha",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"token":      token,
		"image_url":  fmt.Sprintf("/captcha/image/%s", token),
		"expires_in": config.AppConfig.Security.CaptchaTTL,
	})
}

// Image returns the PNG bytes for the given captcha token. Returns 404 if the
// captcha has expired or was never issued.
func (cc *CaptchaController) Image(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Missing captcha token",
		})
		return
	}

	pngBytes, err := cc.captchaService.Render(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Captcha not found or expired",
		})
		return
	}

	c.Data(http.StatusOK, "image/png", pngBytes)
}
