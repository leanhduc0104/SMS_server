package handlers

import (
	"net/http"
	"vcs_server/database"
	"vcs_server/entity"
	"vcs_server/helper"

	"github.com/gin-gonic/gin"
)

type AddUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// Add User godoc
// @Summary Add user
// @Description Add user
//
//	@Tags			User CRUD
//
// @Accept json
// @Produce json
//
//	@Param			user	body		AddUserInput	true	"Add user"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/user [post]
func AddUser(ctx *gin.Context) {
	var input AddUserInput
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot create new user with this information", err.Error(), helper.EmptyObj{}))
		return
	}
	user := entity.User{
		Username: input.Username,
		Role:     input.Role,
	}
	if err := user.HashPassword(input.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	err := database.DB.AddUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Cannot create new user with this information", err.Error(), helper.EmptyObj{}))
		return
	}
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Created successfully new user!", helper.EmptyObj{}))
}
