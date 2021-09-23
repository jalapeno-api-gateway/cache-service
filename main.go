package main

import (
	"log"

	"github.com/jalapeno-api-gateway/cache-service/events"
	"github.com/jalapeno-api-gateway/cache-service/kafka"
	"github.com/jalapeno-api-gateway/cache-service/redis"
)

// Events
var lsNodeEventsChannel = make(chan KafkaEventMessage)
var lsLinkEventsChannel = make(chan KafkaEventMessage)


type KafkaEventMessage struct {
	TopicType	int 	`json:"TopicType,omitempty"`
	Key			string	`json:"_key,omitempty"`
	Id			string	`json:"_id,omitempty"`
	Action		string	`json:"action,omitempty"`
}

func main() {
	log.Print("Starting Cache Service ...")
	redis.InitializeRedisClient()
	kafka.StartEventConsumption()
	redis.InitializeCache()
	events.StartEventProcessing()
}


