package arangodb

import (
	"log"

	"github.com/arangodb/go-driver"
)


func readDocument(document driver.DocumentMeta, err error) bool {
	if driver.IsNoMoreDocuments(err) {
		return false
	}
	if err != nil {
		log.Fatalf("Error while reading from ArangoDb, %v", err)
	}
	return true
}
