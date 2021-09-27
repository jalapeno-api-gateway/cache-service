package main

import (
	"log"
	"os"

	"github.com/jalapeno-api-gateway/arangodb-adapter/arango"
	"github.com/jalapeno-api-gateway/cache-service/events"
	"github.com/jalapeno-api-gateway/cache-service/kafka"
	"github.com/jalapeno-api-gateway/cache-service/redis"
)

func main() {
	log.Print("Starting Cache Service ...")
	arango.InitializeArangoDbAdapter(getDefaultArangoDbConfig())
	redis.InitializeRedisClient()
	kafka.StartEventConsumption()
	redis.InitializeCache()
	events.StartEventProcessing()
}

func getDefaultArangoDbConfig() arango.ArangoDbConfig {
	return arango.ArangoDbConfig{
		Server: os.Getenv("ARANGO_DB"),
		User: os.Getenv("ARANGO_DB_USER"),
		Password: os.Getenv("ARANGO_DB_PASSWORD"),
		DbName: os.Getenv("ARANGO_DB_NAME"),
	}
}

