package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lnb/HRPAuth-Backend-Go/database"
	"github.com/lnb/HRPAuth-Backend-Go/models"
	"github.com/lnb/HRPAuth-Backend-Go/services"
)

type EmailVerificationController struct {
	emailService *services.EmailService
	codeStore    *services.VerificationCodeStore
}

func NewEmailVerificationController() *EmailVerificationController {
	return &EmailVerificationController{
		emailService: services.NewEmailService(),
		codeStore:    services.NewVerificationCodeStore(),
	}
}

type EmailVerificationRequest struct {
	Action  string `json:"action"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Email   string `json:"email"`
	Code    string `json:"code"`
}

func (evc *EmailVerificationController) Handle(c *gin.Context) {
	var req EmailVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	switch req.Action {
	case "send-test-email":
		evc.sendTestEmail(c, req)
	case "send-verification-code":
		evc.sendVerificationCode(c, req)
	case "verify-code":
		evc.verifyCode(c, req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid action",
		})
	}
}

func (evc *EmailVerificationController) sendTestEmail(c *gin.Context, req EmailVerificationRequest) {
	to := req.To
	subject := req.Subject
	message := req.Message

	if to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Recipient email cannot be empty",
		})
		return
	}

	if !isValidEmail(to) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid recipient email format",
		})
		return
	}

	if subject == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email subject cannot be empty",
		})
		return
	}

	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email content cannot be empty",
		})
		return
	}

	err := evc.emailService.SendMail(to, subject, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email sent successfully",
		"data": gin.H{
			"to":      to,
			"subject": subject,
		},
	})
}

func (evc *EmailVerificationController) sendVerificationCode(c *gin.Context, req EmailVerificationRequest) {
	email := req.Email

	if !isValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid email",
		})
		return
	}

	existingCode, found := evc.codeStore.Get(email)
	if found && existingCode != "" {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success": false,
			"message": "Verification code already sent, please wait",
		})
		return
	}

	code := evc.codeStore.GenerateCode()

	if !evc.codeStore.Store(email, code) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to store verification code",
		})
		return
	}

	subject := "HRPAuth - Email Verification Code"
	message := "Your verification code is: " + code + "\n\nThe code is valid for 10 minutes. Please complete the verification as soon as possible.\n\nIf you did not request this code, please ignore this email."

	err := evc.emailService.SendMail(email, subject, message)
	if err != nil {
		evc.codeStore.Delete(email)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification code sent successfully",
	})
}

func (evc *EmailVerificationController) verifyCode(c *gin.Context, req EmailVerificationRequest) {
	email := req.Email
	code := req.Code

	if !isValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid email",
		})
		return
	}

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Verification code is required",
		})
		return
	}

	storedCode, found := evc.codeStore.Get(email)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Verification code expired or not found",
		})
		return
	}

	if code != storedCode {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid verification code",
		})
		return
	}

	evc.codeStore.Delete(email)

	result := database.DB.Model(&models.User{}).
		Where("email = ?", email).
		Update("verified", true)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update verification status",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found or already verified",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification successful",
	})
}
