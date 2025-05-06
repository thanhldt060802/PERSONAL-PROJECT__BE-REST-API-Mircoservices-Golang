package infrastructure

import (
	"fmt"
	"log"
	"thanhldt060802/config"

	"github.com/elastic/go-elasticsearch/v8"
)

var ElasticsearchClient *elasticsearch.Client

func InitElasticsearchClient() {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", config.AppConfig.ElasticsearchHost, config.AppConfig.ElasticsearchPort),
		},
		Username: config.AppConfig.ElasticsearchUsername,
		Password: config.AppConfig.ElasticsearchPassword,
	})
	if err != nil {
		log.Fatal("Connect to Elasticsearch failed: ", err)
	}

	ElasticsearchClient = esClient

	res, err := ElasticsearchClient.Info()
	if err != nil {
		log.Fatal("Ping to Elasticsearch failed: ", err)
	}
	defer res.Body.Close()
	log.Println("Connected to Elasticsearch successful")
}
