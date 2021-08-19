package main

import (
	"log"

	"gitlab.ost.ch/ins/jalapeno-api/graph-db-feeder/events"
	"gitlab.ost.ch/ins/jalapeno-api/graph-db-feeder/kafka"
	"gitlab.ost.ch/ins/jalapeno-api/graph-db-feeder/redis"
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
	log.Print("Starting GraphDbFeeder ...")
	redis.InitializeRedisClient()
	kafka.StartEventConsumption()
	redis.InitializeCache()
	events.StartEventProcessing()
}


