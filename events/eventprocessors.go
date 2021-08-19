package events

import (
	"context"
	"log"

	"gitlab.ost.ch/ins/jalapeno-api/cache-service/arangodb"
	"gitlab.ost.ch/ins/jalapeno-api/cache-service/kafka"
	"gitlab.ost.ch/ins/jalapeno-api/cache-service/redis"
)

func StartEventProcessing() {
	for {
		select {
		case event := <-kafka.LsNodeEvents:
			handleLsNodeEvent(event)
		case event := <-kafka.LsLinkEvents:
			handleLsLinkEvent(event)
		}
	}
}

func handleLsNodeEvent(event kafka.KafkaEventMessage) {
	ctx := context.Background()
	log.Printf("LsNode [%s]: %s\n", event.Action, event.Key)
	if (event.Action == "del") {
		redis.DeleteKey(ctx, event.Key)
	} else {
		updatedDocument := arangodb.FetchLsNode(ctx, event.Key)
		redis.CacheLsNode(updatedDocument.Id, redis.ConvertToRedisLsNode(updatedDocument))
	}
}

func handleLsLinkEvent(event kafka.KafkaEventMessage) {
	ctx := context.Background()
	if (event.Action == "del") {
		redis.DeleteKey(ctx, event.Key)
		log.Printf("LsLink [%s]: %s\n", event.Action, event.Key)
	} else {
		updatedDocument := arangodb.FetchLsLink(ctx, event.Key)
		log.Printf("LsLink [%s]: IGP Metric: %d\n", event.Action, updatedDocument.Igp_metric)
		redis.CacheLsLink(updatedDocument.Id, redis.ConvertToRedisLsLink(updatedDocument))
	}
}