package handlers

import (
	"net/http"
	"vcs_server/entity"
	"vcs_server/helper"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// Import Servers godoc
// @Summary Import servers
// @Description Import servers from excel file
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
// @Param file formData file true "Excel file"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/servers [post]
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func ImportServers(ctx *gin.Context) {
	upload_file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot import servers", err.Error(), helper.EmptyObj{}))
		return
	}

	servers_file, err := excelize.OpenReader(upload_file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot import servers", err.Error(), helper.EmptyObj{}))
		return
	}
	var created_successful_servers []entity.Server
	var created_unsuccessful_servers []entity.Server
	// Get all the rows in the Sheet1.
	rows, err := servers_file.GetRows("Sheet1")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot import servers", err.Error(), helper.EmptyObj{}))
		return
	}
	//skip_row := true
	// for _, row := range rows {
	// 	if skip_row == true {
	// 		skip_row = false
	// 		continue
	// 	}
	// 	new_server := entity.Server{}
	// 	count := 1
	// 	for _, colCell := range row {
	// 		switch count {
	// 		case 1:
	// 			new_server.Name = colCell
	// 		case 2:
	// 			new_server.Ipv4 = colCell
	// 		case 3:
	// 			new_server.Status = colCell
	// 		}
	// 		count += 1
	// 	}
	// 	log.Println(new_server)
	// 	existed_ip := serverService.CheckServerExistence(new_server.Ipv4)
	// 	existed_name := serverService.CheckServerExistence(new_server.Name)
	// 	if existed_ip || existed_name {
	// 		created_unsuccessful_servers = append(created_unsuccessful_servers, new_server)
	// 		continue
	// 	} else {
	// 		err = serverController.CreateServer(&new_server)
	// 		if err != nil {
	// 			created_unsuccessful_servers = append(created_unsuccessful_servers, new_server)
	// 			continue
	// 		}
	// 		created_successful_servers = append(created_successful_servers, new_server)
	// 	}
	// }
	for id, row := range rows {
		if id == 0 {
			continue
		}
		new_server := entity.Server{}
		new_server.Name = row[0]
		new_server.Ipv4 = row[1]
		new_server.Status = row[2]
		err := serverController.CreateServer(&new_server)
		if err != nil {
			created_unsuccessful_servers = append(created_unsuccessful_servers, new_server)
		} else {
			created_successful_servers = append(created_successful_servers, new_server)
		}
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Import task is finished", gin.H{"created_servers": created_successful_servers, "uncreated_servers": created_unsuccessful_servers}))
}
