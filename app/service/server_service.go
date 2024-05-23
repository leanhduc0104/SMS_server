package service

import (
	"fmt"
	"vcs_server/database"
	"vcs_server/entity"
	app_kafka "vcs_server/kafka"
)

type ServerService interface {
	CreateServer(server *entity.Server) error
	ViewServer(id int) (entity.Server, error)
	ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error)
	UpdateServer(server *entity.Server) error
	DeleteServer(id int) error
	CheckServerExistence(ip string) bool
	CheckServerName(name string) bool
}

type serverService struct {
}

func NewServerService() ServerService {
	return &serverService{}
}

func (service *serverService) CreateServer(server *entity.Server) error {
	err := database.DB.CreateServer(server)
	if err != nil {

		return err
	}
	return nil
}

func (service *serverService) ViewServer(id int) (entity.Server, error) {
	server, err := database.DB.ViewServer(id)
	return server, err
}

func (service *serverService) ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error) {
	return database.DB.ViewServers(from, to, perpage, sortby, order, filter)
}

func (service *serverService) UpdateServer(server *entity.Server) error {
	err := database.DB.UpdateServer(server)
	if err != nil {
		return err
	}
	go app_kafka.Producer.ProduceMessage(fmt.Sprint(server.Id))
	return nil
}

func (service *serverService) DeleteServer(id int) error {
	err := database.DB.DeleteServer(id)
	if err != nil {
		return err
	}
	go app_kafka.Producer.ProduceMessage(fmt.Sprint(id))
	return nil
}

func (service *serverService) CheckServerExistence(ip string) bool {
	return database.DB.CheckServerExistence(ip)
}

func (service *serverService) CheckServerName(name string) bool {
	return database.DB.CheckServerName(name)
}
