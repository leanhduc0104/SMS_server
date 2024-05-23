package handlers

import (
	"net/http"
	"strconv"
	"vcs_server/entity"
	"vcs_server/helper"

	"github.com/gin-gonic/gin"
)

// Get Server By Id godoc
// @Summary Get server by id
// @Description Get server by id
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
// @Param id path int true "Server ID"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/server/{id} [get]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func GetServerById(ctx *gin.Context) {
	var server entity.Server
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Server doesn't exist", err.Error(), helper.EmptyObj{}))
		return
	}

	server, err = serverController.ViewServer(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Server doesn't exist", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Server exist", server))
}
