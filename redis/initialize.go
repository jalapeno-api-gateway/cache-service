package redis

import (
	"context"

	"github.com/jalapeno-api-gateway/arangodb-adapter/arango"
	"github.com/jalapeno-api-gateway/model"
)

func InitializeCache() {
	loadLSNodeCollection()
	loadLSLinkCollection()
	loadLSPrefixCollection()
	loadLSSRv6SIDCollection()
}

func loadLSNodeCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLSNodes(ctx)
	for _, document := range documents {
		CacheObject(document.ID, model.ConvertLSNode(document))
	}
}

func loadLSLinkCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLSLinks(ctx)
	for _, document := range documents {
		CacheObject(document.ID, model.ConvertLSLink(document))
	}
}

func loadLSPrefixCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLSPrefix(ctx)
	for _, document := range documents {
		CacheObject(document.ID, model.ConvertLSPrefix(document))
	}
}

func loadLSSRv6SIDCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLSSRv6SID(ctx)
	for _, document := range documents {
		CacheObject(document.ID, model.ConvertLSSRv6SID(document))
	}
}