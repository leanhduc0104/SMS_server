package handlers

import (
	"net/http"
	"os"
	"time"
	"vcs_server/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Login
// @Description Login
//
//	@Tags			Auth
//
// @Accept json
// @Produce json
//
//	@Param			login	body		LoginInput	true	"Login"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /login [post]
func Login(c *gin.Context) {
	jwtSecret := os.Getenv("JWT_SECRET")
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := database.DB.FindUserByUsername(input.Username)

	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
