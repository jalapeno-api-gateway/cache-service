package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	graphproto "gitlab.ost.ch/ins/jalapeno-api/graph-db-feeder/proto"
	grpc "google.golang.org/grpc"
)

type graphDbFeederService struct {
	graphproto.UnimplementedGraphDbFeederServer
}

type nodeDocument struct {
	RouterID string
	Name string;
	ASN int32;
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
	//log.Print("GraphDbFeeder Up and Running")
}

func (s *graphDbFeederService) GetNodes(nodeIds *graphproto.NodeIds, responseStream graphproto.GraphDbFeeder_GetNodesServer) error {
	//Access ArangoDB and get Nodes collection
	arangoDbClient := connectToArangoDb()
	nodes := getNodesFromArangoDb(arangoDbClient, nodeIds.Ids)
	//Return stream of Nodes
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
		Connection: conn,
		Authentication: driver.BasicAuthentication(os.Getenv("ARANGO_DB_USER"), os.Getenv("ARANGO_DB_PASSWORD")),
	})
	if err != nil {
		log.Fatalf("Could not create new ArangoDb client, %v", err)
	}
	return c
}

func getNodesFromArangoDb(arangoDbClient driver.Client, ids []int32) []graphproto.Node {
	var nodes []graphproto.Node
	
	// Opening the database
	ctx := context.Background()
	db, err := arangoDbClient.Database(ctx, os.Getenv("ARANGO_DB_NAME"))
	if err != nil {
		log.Fatalf("Could not open database, %v", err)
	}

	// Opening LSNode collection
	col, err := db.Collection(ctx, "LSNode")
	if err != nil {
		log.Fatalf("Could not open LSNode collection, %v", err)
	}

	for _, id := range ids {
		_id := getDbId(id)
		
		// Reading document from collection
		var doc nodeDocument
		_, err := col.ReadDocument(ctx, _id, &doc)
		if err != nil {
			log.Fatalf("Could not read document with _id: %s, %v", _id, err)
		}

		node := graphproto.Node{Id: getIdFromRouterId(doc.RouterID), Name: doc.Name, Asn: doc.ASN}
		nodes = append(nodes, node)
	}

	return nodes
}

func getDbId(id int32) string {
	// Convert int id = 1 to string _id = 1.1.1.1
	return fmt.Sprintf("%d.%d.%d.%d", id, id, id, id)
}

func getIdFromRouterId(_id string) int32 {
	// Convert string _id = 1.1.1.1 to int id = 1
	re := regexp.MustCompile("[0-9]+")
	id, err := strconv.Atoi(re.FindAllString(_id, -1)[0])
	if err != nil {
		log.Fatalf("Could not extract id from string, %v", err)
	}
	return int32(id)
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
