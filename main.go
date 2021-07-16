package main

import (
	"context"
	"log"
	"net"
	"os"
	"runtime"

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
	//Sequentiell
	//arangoDbClient := connectToArangoDb()
	// nodes := getNodesFromArangoDb(arangoDbClient, nodeIds.Ids) //sequentiell sollte concurrent sein.
	// for _, node := range nodes {
	// 	if err := responseStream.Send(&node); err != nil {
	// 		log.Fatalf("Could not return node to request-service, %v", err)
	// 	}
	// }
	// return nil

	//Concurrent try 1 with worker pool
	log.Print("GetNodes called from RequestService")
	log.Print("Start fetching Nodes")
	var workerId = 1
	jobs := make(chan []string, len(nodeIds.Ids))             //jobs contains ids which need to be fetched from DB; Buffer Size = amount of ids to fetch
	results := make(chan []graphproto.Node, len(nodeIds.Ids)) //results contains fetched node objects
	for i := 0; i < runtime.NumCPU(); i++ {                   //create as many workers as cores exist
		go worker(jobs, results, workerId) //start worker to fetch nodes from DB
		workerId++
	}
	for _, id := range nodeIds.Ids {
		jobs <- []string{id} //fill jobs queue: only one id per array can/should be adjusted
	}
	close(jobs)                             //all jobs created so channel can be closed
	for j := 0; j < len(nodeIds.Ids); j++ { //for each job one result is expected
		nodes := <-results
		for _, node := range nodes {
			if err := responseStream.Send(&node); err != nil {
				log.Fatalf("Could not return node to request-service, %v", err)
			}
		}
	}
	log.Print("Finished Fetching Nodes")
	return nil
}

func worker(jobs <-chan []string, results chan<- []graphproto.Node, workerId int) {
	arangoDbClient := connectToArangoDb()
	log.Printf("Worker %d fetching from DB", workerId)
	for job := range jobs {
		results <- getNodesFromArangoDb(arangoDbClient, job)
	}
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
