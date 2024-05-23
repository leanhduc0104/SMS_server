package handlers

import (
	"net/http"
	"strconv"
	"vcs_server/helper"
	"vcs_server/mail"

	"github.com/gin-gonic/gin"
)

// Report Server Information Intentionally godoc
// @Summary Report server information intentionally
// @Description Report server information
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
// @Param time query int false "Time"
// @Success 200 {object} helper.Response
// @Router /api/servers/report [get]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func ReportServer(ctx *gin.Context) {
	timeQuery := ctx.DefaultQuery("time", "12")
	timeInt, err := strconv.Atoi(timeQuery)
	if err != nil {
		// Nếu có lỗi trong quá trình chuyển đổi, trả về lỗi bad request
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid time query parameter",
		})
		return
	}

	err = mail.SendReport(timeInt)
	if err != nil {
		// Nếu có lỗi trong quá trình gửi mail, trả về lỗi internal server error
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot send mail", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Sending server report successfully", helper.EmptyObj{}))
}
