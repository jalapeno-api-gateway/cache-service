package redis

import (
	"context"

	"github.com/jalapeno-api-gateway/arangodb-adapter/arango"
	"github.com/jalapeno-api-gateway/model/topology"
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
		CacheObject(document.ID, topology.ConvertLSNode(document))
	}
}

func loadLSLinkCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLSLinks(ctx)
	for _, document := range documents {
		CacheObject(document.ID, topology.ConvertLSLink(document))
	}
}

func loadLSPrefixCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLSPrefix(ctx)
	for _, document := range documents {
		CacheObject(document.ID, topology.ConvertLSPrefix(document))
	}
}

func loadLSSRv6SIDCollection() {
	ctx := context.Background()
	documents := arango.FetchAllLSSRv6SID(ctx)
	for _, document := range documents {
		CacheObject(document.ID, topology.ConvertLSSRv6SID(document))
	}
}