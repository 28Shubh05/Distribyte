package handlers

import (
	"Distribyte/backend/services"
	"Distribyte/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request",
		})

		return
	}

	user, err := services.LoginUser(
		req.Email,
		req.Password,
	)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid credentials",
		})

		return
	}

	token, err :=
		utils.GenerateToken(
			user.ID,
			user.Email,
		)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Token generation failed",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}

func Register(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request",
		})

		return
	}

	err := services.RegisterUser(
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User registered",
	})
}
