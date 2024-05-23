package healcheck

import (
	"context"
	"log"
	"os"
	"time"
	"vcs_server/database"
	"vcs_server/entity"
	"vcs_server/helper"
	"vcs_server/service"

	"github.com/olivere/elastic/v7"
)

func logToElasticsearch(client *elastic.Client, index string, server entity.Server) error {
	doc := map[string]interface{}{
		"ipv4": server.Ipv4,
		"time": time.Now(),
	}

	_, err := client.Index().
		Index(index).
		BodyJson(doc).
		Do(context.Background())
	return err
}

func worker(index string, servers []entity.Server, client *elastic.Client, server_service service.ServerService) {
	for _, server := range servers {
		go func(server entity.Server) {
			if err := helper.PingServer(server.Name); err != nil {
				if err := logToElasticsearch(client, index, server); err != nil {
					log.Printf("Failed to log to Elasticsearch: %v", err)
				}
				if server.Status == "up" {
					server.Status = "down"
					server_service.UpdateServer(&server)
				}
			} else {
				if server.Status == "down" {
					server.Status = "up"
					server_service.UpdateServer(&server)
				}
			}
		}(server)
	}
}

func SendHealcheck() {
	index := os.Getenv("ELASTICSEARCH_INDEX")
	client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ELASTICSEARCH_URL")), elastic.SetSniff(false))

	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to check if index exists: %v", err)
	}

	if !exists {
		_, err := client.CreateIndex(index).Do(context.Background())
		if err != nil {
			log.Fatalf("Failed to create index: %v", err)
		}
	}

	server_service := service.NewServerService()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		servers, err := database.DB.ViewServers(0, 0, 0, "", "", "")
		if err != nil {
			log.Printf("Failed to get servers from DB: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		worker(index, servers, client, server_service)

		<-ticker.C
	}
}
