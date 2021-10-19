package events

import (
	"context"

	"github.com/jalapeno-api-gateway/jagw-core/arango"
	"github.com/jalapeno-api-gateway/cache-service/kafka"
	"github.com/jalapeno-api-gateway/cache-service/redis"
	"github.com/jalapeno-api-gateway/jagw-core/model/topology"
	"github.com/jalapeno-api-gateway/jagw-core/model/class"
)

func StartEventProcessing() {
	for {
		select {
			case event := <-kafka.LsNodeEvents: handleEvent(event, class.LsNode)
			case event := <-kafka.LsLinkEvents: handleEvent(event, class.LsLink)
			case event := <-kafka.LsPrefixEvents: handleEvent(event, class.LsPrefix)
			case event := <-kafka.LsSrv6SidEvents: handleEvent(event, class.LsSrv6Sid)
		}
	}
}

func handleEvent(event kafka.KafkaEventMessage, className class.Class) {
	ctx := context.Background()
	if (event.Action == "del") {
		redis.DeleteKey(ctx, event.Key)
	} else {
		id, obj := fetchDocument(ctx, event.Key, className)
		redis.CacheObject(id, obj)
	}
}

func fetchDocument(ctx context.Context, key string, className class.Class) (string, interface{}) {
	switch className {
		case class.LsNode:
			doc := arango.FetchLsNode(ctx, key)
			return doc.ID, topology.ConvertLsNode(doc)
		case class.LsLink:
			doc := arango.FetchLsLink(ctx, key)
			return doc.ID, topology.ConvertLsLink(doc)
		case class.LsPrefix:
			doc := arango.FetchLsPrefix(ctx, key)
			return doc.ID, topology.ConvertLsPrefix(doc)
		case class.LsSrv6Sid:
			doc := arango.FetchLsSrv6Sid(ctx, key)
			return doc.ID, topology.ConvertLsSrv6Sid(doc)
		default: return "", nil
	}
}