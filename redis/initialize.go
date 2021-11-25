package redis

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jalapeno-api-gateway/jagw-core/arango"
	"github.com/jalapeno-api-gateway/jagw-core/model/topology"
)

func InitializeRedisClient() {
	redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    os.Getenv("SENTINEL_MASTER"),
		SentinelAddrs: []string{os.Getenv("SENTINEL_ADDRESS")},
		Password:      os.Getenv("REDIS_PASSWORD"),
		DB:            0,
	})
}

func InitializeCache() {
	loadLsNodeCollection()
	loadLsLinkCollection()
	loadLsPrefixCollection()
	loadLsSrv6SidCollection()
	loadLsNodeEdgeCollection()
	loadLsNodeCoordinatesCollection()
}

func loadLsNodeCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLsNodes(ctx)
	for _, document := range documents {
		CacheObject(document.ID, topology.ConvertLsNode(document))
	}
}

func loadLsLinkCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLsLinks(ctx)
	for _, document := range documents {
		CacheObject(document.ID, topology.ConvertLsLink(document))
	}
}

func loadLsPrefixCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLsPrefixes(ctx)
	for _, document := range documents {
		CacheObject(document.ID, topology.ConvertLsPrefix(document))
	}
}

func loadLsSrv6SidCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLsSrv6Sids(ctx)
	for _, document := range documents {
		CacheObject(document.ID, topology.ConvertLsSrv6Sid(document))
	}
}

func loadLsNodeEdgeCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLsNodeEdges(ctx)
	for _, document := range documents {
		CacheObject(document.ID, topology.ConvertLsNodeEdge(document))
	}
}

func loadLsNodeCoordinatesCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLsNodeCoordinates(ctx)
	for _, document := range documents {
		log.Printf("%s\n", document.Key)
		CacheObject(document.ID, topology.ConvertLsNodeCoordinates(document))
	}
}