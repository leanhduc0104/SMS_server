package helper

import (
	"vcs_server/database"
	"vcs_server/entity"
)

func Migration() {
	type AddUserInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	user_input := AddUserInput{Username: "duc", Password: "123456", Role: "admin"}
	user := entity.User{Username: user_input.Username, Role: user_input.Role}
	if err := user.HashPassword(user_input.Password); err != nil {
		return
	}
	database.DB.AddUser(&user)
}
