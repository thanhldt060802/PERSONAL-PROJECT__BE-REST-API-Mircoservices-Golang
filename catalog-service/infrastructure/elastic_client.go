package infrastructure

import (
	"fmt"
	"log"
	"thanhldt060802/config"

	"github.com/elastic/go-elasticsearch/v8"
)

var ESClient *elasticsearch.Client

func InitElasticsearchClient() {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", config.AppConfig.ESHost, config.AppConfig.ESPort),
		},
		Username: config.AppConfig.ESUsername,
		Password: config.AppConfig.ESPassword,
	})
	if err != nil {
		log.Fatal("Connect to ES failed: ", err)
	}

	ESClient = esClient

	res, err := ESClient.Info()
	if err != nil {
		log.Fatal("Ping to ES failed: ", err)
	}
	defer res.Body.Close()
	log.Println("Connected to ES successful")

}
