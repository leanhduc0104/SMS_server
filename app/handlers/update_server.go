package handlers

import (
	"net/http"
	"strconv"
	"vcs_server/entity"
	"vcs_server/helper"

	"github.com/gin-gonic/gin"
)

// Update Server godoc
// @Summary Update server
// @Description Update server by id
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
// @Param id path int true "Server ID"
// @Param server body entity.Server true "Update server"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/server/{id} [put]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func UpdateServer(ctx *gin.Context) {
	var server entity.Server
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot update server with this information", err.Error(), helper.EmptyObj{}))
		return
	}

	if err := ctx.ShouldBind(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot update server with this information", err.Error(), helper.EmptyObj{}))
		return
	}

	_ = id

	server.Id = id

	err = serverController.UpdateServer(&server)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot update server with this information", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Updated server successfully!", helper.EmptyObj{}))

}
