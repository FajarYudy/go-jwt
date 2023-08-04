package controllers

import (
	"go-jwt/auth"
	"go-jwt/db"
	"go-jwt/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Username    string `json:"username"`
	Type 		string `json:"type"`
	Token 		string `json:"token"`	
}

func GenerateToken(c *gin.Context) {
	var request TokenRequest
	var user models.User
		
	if err := c.ShouldBindJSON(&request); err != nil {
		PayloadValidation(c, err)		 
		return
	}

	// check if email exists and password is correct
	record := db.Gorm.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {		
		PayloadError(c, http.StatusInternalServerError, record.Error.Error())
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {		
		PayloadError(c, http.StatusInternalServerError, "invalid credentials")
		return
	}

	tokenString, err:= auth.GenerateJWT(user.Email, user.Username)
	if err != nil {		
		PayloadError(c, http.StatusInternalServerError, err.Error())
		return
	}	

	Payload(c, TokenResponse{Username: user.Username, Type: "Bearer", Token: tokenString})
}