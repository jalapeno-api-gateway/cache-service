package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jalapeno-api-gateway/cache-service/events"
	"github.com/jalapeno-api-gateway/cache-service/kafka"
	"github.com/jalapeno-api-gateway/cache-service/redis"
	"github.com/jalapeno-api-gateway/jagw-core/arango"
)

func main() {
	log.Print("Starting Cache Service ...")
	arango.InitializeArangoDbAdapter(getDefaultArangoDbConfig())
	log.Print("1")
	redis.InitializeRedisClient()
	log.Print("2")
	kafka.StartEventConsumption()
	log.Print("3")
	redis.InitializeCache()
	log.Print("4")
	events.StartEventProcessing()
	log.Print("5")
}

func getDefaultArangoDbConfig() arango.ArangoDbConfig {
	return arango.ArangoDbConfig{
		Server: fmt.Sprintf("http://%s", os.Getenv("ARANGO_ADDRESS")),
		User: os.Getenv("ARANGO_DB_USER"),
		Password: os.Getenv("ARANGO_DB_PASSWORD"),
		DbName: os.Getenv("ARANGO_DB_NAME"),
	}
}

