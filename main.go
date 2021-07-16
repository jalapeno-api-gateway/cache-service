package main

import (
	"context"
	"log"
	"net"
	"os"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	graphproto "gitlab.ost.ch/ins/jalapeno-api/graph-db-feeder/proto"
	grpc "google.golang.org/grpc"
)

type graphDbFeederService struct {
	graphproto.UnimplementedGraphDbFeederServer
}

type NodeDocument struct {
	Key       string `json:"_key,omitempty"`
	Name      string `json:"name,omitempty"`
	Asn       int32  `json:"asn,omitempty"`
	Router_ip string `json:"router_ip,omitempty"`
}

func newServer() *graphDbFeederService {
	s := &graphDbFeederService{}
	return s
}

func main() {
	//Start gRPC server for Request Service
	log.Print("Starting GraphDbFeeder ...")
	lis, err := net.Listen("tcp", "0.0.0.0:9001")
	if err != nil {
		log.Fatalf("Failed to listen on port 9001: %v", err)
	}
	grpcServer := grpc.NewServer()
	graphproto.RegisterGraphDbFeederServer(grpcServer, newServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9001: %v", err)
	}
}

func (s *graphDbFeederService) GetNodes(nodeIds *graphproto.NodeIds, responseStream graphproto.GraphDbFeeder_GetNodesServer) error {
	arangoDbClient := connectToArangoDb()
	nodes := getNodesFromArangoDb(arangoDbClient, nodeIds.Ids)
	for _, node := range nodes {
		if err := responseStream.Send(&node); err != nil {
			log.Fatalf("Could not return node to request-service, %v", err)
		}
	}
	return nil
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

func getNodesFromArangoDb(arangoDbClient driver.Client, keys []string) []graphproto.Node {
	var nodes []graphproto.Node

	ctx := context.Background()
	db, err := arangoDbClient.Database(ctx, os.Getenv("ARANGO_DB_NAME"))
	if err != nil {
		log.Fatalf("Could not open database, %v", err)
	}

	col, err := db.Collection(ctx, "LSNode")
	if err != nil {
		log.Fatalf("Could not open LSNode collection, %v", err)
	}

	log.Println("1")
	for _, key := range keys {
		var doc NodeDocument
		_, err := col.ReadDocument(ctx, key, &doc)
		if err != nil {
			log.Fatalf("Could not read document with _id: %s, %v", key, err)
		}
		node := graphproto.Node{Key: doc.Key, Name: doc.Name, Asn: doc.Asn, RouterIp: doc.Router_ip}
		nodes = append(nodes, node)
	}
	return nodes
}

// func cacheNode(ctx context.Context, nodeId int, node *rs.NodeResponse) {
// 	key := strconv.Itoa(nodeId)
// 	redis.StoreMessage(ctx, key, node)
// }

// func getNodeFromCache(ctx context.Context, nodeId int) *rs.NodeResponse {
// 	key := strconv.Itoa(nodeId)
// 	return redis.ReadMessage(ctx, key)
// }

// func processGetNodeRequest(id int) rs.NodeResponse {
// 	return rs.NodeResponse{
// 		Id:   int32(id),
// 		Name: getNodeNameById(id),
// 	}
// }

// func getNodeNameById(id int) string {
// 	return "Node-" + strconv.Itoa(id)
// }
