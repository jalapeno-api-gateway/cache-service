package events

import (
	"context"

	"github.com/jalapeno-api-gateway/arangodb-adapter/arango"
	"github.com/jalapeno-api-gateway/cache-service/kafka"
	"github.com/jalapeno-api-gateway/cache-service/redis"
	"github.com/jalapeno-api-gateway/model/topology"
)

type EventType string

const (
	LSNodeEvent EventType = "LSNodeEvent"
	LSLinkEvent EventType = "LSLinkEvent"
	LSPrefixEvent EventType = "LSPrefixEvent"
	LSSIDv6SIDEvent EventType = "LSSIDv6SIDEvent"
	PhysicalInterfaceTelemetryEvent EventType = "PhysicalInterfaceTelemetryEvent"
	LoopbackInterfaceTelemetryEvent EventType = "LoopbackInterfaceTelemetryEvent"
)

func StartEventProcessing() {
	for {
		select {
			case event := <-kafka.LSNodeEvents: handleEvent(event, LSNodeEvent)
			case event := <-kafka.LSLinkEvents: handleEvent(event, LSLinkEvent)
			case event := <-kafka.LSPrefixEvents: handleEvent(event, LSPrefixEvent)
			case event := <-kafka.LSSRV6SIDEvents: handleEvent(event, LSSIDv6SIDEvent)
		}
	}
}

func handleEvent(event kafka.KafkaEventMessage, eventType EventType) {
	ctx := context.Background()
	if (event.Action == "del") {
		redis.DeleteKey(ctx, event.Key)
	} else {
		id, obj := fetchDocument(ctx, event.Key, eventType)
		redis.CacheObject(id, obj)
	}
}

func fetchDocument(ctx context.Context, key string, eventType EventType) (string, interface{}) {
	switch eventType {
		case LSNodeEvent:
			doc := arango.FetchLSNode(ctx, key)
			return doc.ID, topology.ConvertLSNode(doc)
		case LSLinkEvent:
			doc := arango.FetchLSLink(ctx, key)
			return doc.ID, topology.ConvertLSLink(doc)
		case LSPrefixEvent:
			doc := arango.FetchLSPrefix(ctx, key)
			return doc.ID, topology.ConvertLSPrefix(doc)
		case LSSIDv6SIDEvent:
			doc := arango.FetchLSSRv6SID(ctx, key)
			return doc.ID, topology.ConvertLSSRv6SID(doc)
		default: return "", nil
	}
}