package arangodb

import (
	"context"
	"log"
	"os"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)


func queryArangoDbDatabase(ctx context.Context, query string) driver.Cursor {
	db := openArangoDbDatabase(ctx)
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		log.Fatalf("Could not create Cursor , %v", err)
	}
	defer cursor.Close()
	return cursor
}

func openArangoDbDatabase(ctx context.Context) driver.Database {
	arangoDbClient := connectToArangoDb()
	db, err := arangoDbClient.Database(ctx, os.Getenv("ARANGO_DB_NAME"))
	if err != nil {
		log.Fatalf("Could not open database, %v", err)
	}
	return db
}

func connectToArangoDb() driver.Client {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{os.Getenv("ARANGO_DB")},
	})
	if err != nil {
		log.Fatalf("Could not connect to ArangoDb, %v", err)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(os.Getenv("ARANGO_DB_USER"), os.Getenv("ARANGO_DB_PASSWORD")),
	})
	if err != nil {
		log.Fatalf("Could not create new ArangoDb client, %v", err)
	}
	return c
}