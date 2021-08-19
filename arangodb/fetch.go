package arangodb

import (
	"context"
)

//
// ---> FETCH SINGLE <---
//

func FetchLsNode(ctx context.Context, key string) LsNodeDocument {
	cursor := queryArangoDbDatabase(ctx, "FOR d IN LSNode FILTER d._key == \"" + key + "\" RETURN d");
	var document LsNodeDocument
	readDocument(cursor.ReadDocument(ctx, &document))
	return document
}

func FetchLsLink(ctx context.Context, key string) LsLinkDocument {
	cursor := queryArangoDbDatabase(ctx, "FOR d IN LSLink FILTER d._key == \"" + key + "\" RETURN d");
	var document LsLinkDocument
	readDocument(cursor.ReadDocument(ctx, &document))
	return document
}

//
// ---> FETCH ALL <---
//

func FetchAllLsNodes(ctx context.Context) []LsNodeDocument {
	cursor := queryArangoDbDatabase(ctx, "FOR d IN LSNode RETURN d");
	var documents []LsNodeDocument
	for {
		var document LsNodeDocument
		if (!readDocument(cursor.ReadDocument(ctx, &document))) {
			break
		}
		documents = append(documents, document)
	}
	return documents
}

func FetchAllLsLinks(ctx context.Context) []LsLinkDocument {
	cursor := queryArangoDbDatabase(ctx, "FOR d IN LSLink RETURN d");
	var documents []LsLinkDocument
	for {
		var document LsLinkDocument
		if (!readDocument(cursor.ReadDocument(ctx, &document))) {
			break
		}
		documents = append(documents, document)
	}
	return documents
}