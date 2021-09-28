package redis

import (
	"context"

	"github.com/jalapeno-api-gateway/arangodb-adapter/arango"
)

func InitializeCache() {
	loadLSNodeCollection()
	loadLSLinkCollection()
}

func loadLSNodeCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLsNodes(ctx)
	for _, document := range documents {
		CacheLsNode(document.Id, ConvertToRedisLsNode(document))
	}
}

func loadLSLinkCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLsLinks(ctx)
	for _, document := range documents {
		CacheLsLink(document.Id, ConvertToRedisLsLink(document))
	}
}