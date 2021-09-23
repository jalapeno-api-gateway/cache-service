package redis

import (
	"context"

	"github.com/jalapeno-api-gateway/cache-service/arangodb"
)

func InitializeCache() {
	loadLSNodeCollection()
	loadLSLinkCollection()
}

func loadLSNodeCollection() {
	ctx := context.Background()
	documents := arangodb.FetchAllLsNodes(ctx)
	for _, document := range documents {
		CacheLsNode(document.Id, ConvertToRedisLsNode(document))
	}
}

func loadLSLinkCollection() {
	ctx := context.Background()
	documents := arangodb.FetchAllLsLinks(ctx)
	for _, document := range documents {
		CacheLsLink(document.Id, ConvertToRedisLsLink(document))
	}
}