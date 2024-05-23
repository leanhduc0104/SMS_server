package handlers

import (
	"net/http"
	"strconv"
	"vcs_server/helper"

	"github.com/gin-gonic/gin"
)

// Delete Server godoc
// @Summary Delete server
// @Description Delete server by id
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
// @Router /api/server/{id} [delete]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func DeleteServer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot delete server with this id", err.Error(), helper.EmptyObj{}))
		return
	}

	err = serverController.DeleteServer(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot delete server with this id", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Deleted server successfully!", helper.EmptyObj{}))

}
