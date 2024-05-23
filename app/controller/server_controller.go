package controller

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
	"vcs_server/cache"
	"vcs_server/entity"
	"vcs_server/helper"
	"vcs_server/service"

	"github.com/olivere/elastic/v7"
)

type ServerController interface {
	CreateServer(server *entity.Server) error
	ViewServer(id int) (entity.Server, error)
	ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error)
	UpdateServer(server *entity.Server) error
	DeleteServer(id int) error
	ReportServerStatus(time_report int) ([]helper.ReportReposne, error)
}

type serverController struct {
	serverService service.ServerService
	serverCache   cache.ServerCache
}

func NewServerController(serverService service.ServerService, serverCache cache.ServerCache) ServerController {
	return &serverController{
		serverService: serverService,
		serverCache:   serverCache,
	}
}

func (controller *serverController) CreateServer(server *entity.Server) error {
	err := controller.serverService.CreateServer(server)
	if err != nil {
		return err
	}
	return nil

}

func (controller *serverController) ViewServer(id int) (entity.Server, error) {
	var server *entity.Server = controller.serverCache.Get(fmt.Sprint(id))
	if server == nil {
		server, err := controller.serverService.ViewServer(id)
		controller.serverCache.Set(fmt.Sprint(id), &server)
		return server, err
	}
	return *server, nil
}

func (controller *serverController) ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error) {
	log.Println("hello")
	key := strconv.FormatInt(int64(from+2), 10) + "-" + strconv.FormatInt(int64(to+2), 10) + "-" + strconv.FormatInt(int64(perpage+2), 10) + "-" + sortby + "-" + order + "-" + filter
	var servers []entity.Server = controller.serverCache.AGet(key)
	if servers == nil {
		result, err := controller.serverService.ViewServers(from, to, perpage, sortby, order, filter)
		controller.serverCache.ASet(key, result)
		return result, err
	}
	return controller.serverService.ViewServers(from, to, perpage, sortby, order, filter)
}

func (controller *serverController) UpdateServer(server *entity.Server) error {
	err := controller.serverService.UpdateServer(server)
	if err != nil {
		return err
	}
	return nil
}

func (controller *serverController) DeleteServer(id int) error {
	err := controller.serverService.DeleteServer(id)
	return err
}

func (controller *serverController) ReportServerStatus(time_report int) ([]helper.ReportReposne, error) {
	client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ELASTICSEARCH_URL")), elastic.SetSniff(false))
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
		return nil, err
	}
	servers, err := controller.serverService.ViewServers(0, 0, 0, "", "", "")
	if err != nil {
		log.Printf("Failed to get servers from DB: %v", err)
		return nil, err
	}
	index := os.Getenv("ELASTICSEARCH_INDEX")
	startTime := time.Now().Add(time.Duration(-time_report) * time.Hour).Format(time.RFC3339)
	endTime := time.Now().Format(time.RFC3339)
	var wg sync.WaitGroup
	var reports []helper.ReportReposne
	wg.Add(len(servers))
	for _, server := range servers {
		go func(server entity.Server) {
			query := elastic.NewBoolQuery().Must(
				elastic.NewTermQuery("ipv4", server.Ipv4),
			).Filter(
				elastic.NewRangeQuery("time").Gte(startTime).Lte(endTime),
			)
			searchResult, err := client.Search().
				Index(index). // Đặt tên index của bạn ở đây
				Query(query).
				Size(0).
				Do(context.Background())
			if err != nil {
				log.Fatalf("Failed to execute search query: %v", err)
			}

			// In ra kết quả
			totalHits := searchResult.Hits.TotalHits.Value

			// Lặp qua kết quả và in ra các bản ghi
			downtime := float64(totalHits) / float64(time_report*12)
			uptime := 100 - math.Round(downtime*10000)/100
			reports = append(reports, helper.ReportReposne{Server: server, Uptime: uptime})
			wg.Done()
		}(server)

	}
	wg.Wait()
	return reports, err
}
