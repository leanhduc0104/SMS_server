package handlers

import (
	"net/http"
	"vcs_server/cache"
	"vcs_server/controller"
	"vcs_server/entity"
	"vcs_server/service"

	"vcs_server/helper"

	"github.com/gin-gonic/gin"
)

var (
	serverCache      cache.ServerCache           = cache.NewRedisCache("redis:6379", 0, 120)
	serverService    service.ServerService       = service.NewServerService()
	serverController controller.ServerController = controller.NewServerController(serverService, serverCache)
)

// Create New Server godoc
// @Summary Create new server
// @Description Create new server with provided information
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
//	@Param			server	body		entity.Server	true	"Add server"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/server/ [post]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func CreateNewServer(ctx *gin.Context) {
	var server entity.Server
	if err := ctx.ShouldBind(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot create new server with this information", err.Error(), helper.EmptyObj{}))
		return
	}
	err := serverController.CreateServer(&server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot create new server with this information", err.Error(), helper.EmptyObj{}))
		return
	}
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Created successfully new server!", helper.EmptyObj{}))
}
