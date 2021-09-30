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
			case event := <-kafka.LSNodeEvents: handleEvent(event, class.LSNode)
			case event := <-kafka.LSLinkEvents: handleEvent(event, class.LSLink)
			case event := <-kafka.LSPrefixEvents: handleEvent(event, class.LSPrefix)
			case event := <-kafka.LSSRV6SIDEvents: handleEvent(event, class.LSSRv6SID)
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
		case class.LSNode:
			doc := arango.FetchLSNode(ctx, key)
			return doc.ID, topology.ConvertLSNode(doc)
		case class.LSLink:
			doc := arango.FetchLSLink(ctx, key)
			return doc.ID, topology.ConvertLSLink(doc)
		case class.LSPrefix:
			doc := arango.FetchLSPrefix(ctx, key)
			return doc.ID, topology.ConvertLSPrefix(doc)
		case class.LSSRv6SID:
			doc := arango.FetchLSSRv6SID(ctx, key)
			return doc.ID, topology.ConvertLSSRv6SID(doc)
		default: return "", nil
	}
}