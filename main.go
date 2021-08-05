package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"runtime"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	redis "github.com/go-redis/redis/v8"
	graphproto "gitlab.ost.ch/ins/jalapeno-api/graph-db-feeder/proto"
	grpc "google.golang.org/grpc"
)

// GLOBAL VARS
var rdb = redis.NewFailoverClient(&redis.FailoverOptions{
	MasterName:    os.Getenv("SENTINEL_MASTER"),
	SentinelAddrs: []string{os.Getenv("SENTINEL_ADDRESS")},
	Password:      os.Getenv("REDIS_PASSWORD"),
	DB:            0,
})

//Type Definition
type graphDbFeederService struct {
	graphproto.UnimplementedGraphDbFeederServer
}

type NodeDocument struct {
	Key       string `json:"_key,omitempty"`
	Name      string `json:"name,omitempty"`
	Asn       int32  `json:"asn,omitempty"`
	Router_ip string `json:"router_ip,omitempty"`
}

//Marshaling Method used to write NodeDocument to Redis
func (node NodeDocument) MarshalBinary() ([]byte, error) {
	return json.Marshal(node)
}

type LinkDocument struct {
	Key           string `json:"_key,omitempty"`
	Router_ip     string `json:"router_ip,omitempty"`
	Peer_ip       string `json:"peer_ip,omitempty"`
	LocalLink_ip  string `json:"local_link_ip,omitempty"`
	RemoteLink_ip string `json:"remote_link_ip,omitempty"`
}

func (link LinkDocument) MarshalBinary() ([]byte, error) {
	return json.Marshal(link)
}

func newServer() *graphDbFeederService {
	s := &graphDbFeederService{}
	return s
}

//Implementd GRPC Methods
func (s *graphDbFeederService) GetNodes(nodeIds *graphproto.NodeIds, responseStream graphproto.GraphDbFeeder_GetNodesServer) error {
	//Concurrent with worker pool
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

func main() {
	log.Print("Starting GraphDbFeeder ...")
	loadArangoDbIntoCache()
	lis, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}
	grpcServer := grpc.NewServer()
	graphproto.RegisterGraphDbFeederServer(grpcServer, newServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}

func loadArangoDbIntoCache() {
	log.Print("Loading LSNode Collection from ArangoDb into Cache")
	loadLSNodeCollection()
	log.Print("Loading of LSNode Collection from ArangoDb into Cache DONE")

	log.Print("Loading LSLink Collection from ArangoDb into Cache")
	loadLSLinkCollection()
	log.Print("Loading LSLink Collection from ArangoDb into Cache DONE")
}

func loadLSNodeCollection() {
	ctx := context.Background()
	arangoDbClient := connectToArangoDb()
	db, err := arangoDbClient.Database(ctx, os.Getenv("ARANGO_DB_NAME"))
	if err != nil {
		log.Fatalf("Could not open database, %v", err)
	}
	cursor, err := db.Query(ctx, "FOR d IN LSNode RETURN d", nil)
	if err != nil {
		log.Fatalf("Could not create Cursor , %v", err)
	}
	defer cursor.Close()
	for {
		var node NodeDocument
		meta, err := cursor.ReadDocument(ctx, &node)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Fatalf("Could not fetch Node from LSNode Collection , %v", err)
		}
		writeNodeToRedis(meta.Key, node)
	}
}

func loadLSLinkCollection() {
	ctx := context.Background()
	arangoDbClient := connectToArangoDb()
	db, err := arangoDbClient.Database(ctx, os.Getenv("ARANGO_DB_NAME"))
	if err != nil {
		log.Fatalf("Could not open database, %v", err)
	}
	cursor, err := db.Query(ctx, "FOR d IN LSLink RETURN d", nil)
	if err != nil {
		log.Fatalf("Could not create Cursor , %v", err)
	}
	defer cursor.Close()
	for {
		var node LinkDocument
		meta, err := cursor.ReadDocument(ctx, &node)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Fatalf("Could not fetch Link from LSLink Collection , %v", err)
		}
		writeLinkToRedis(meta.Key, node)
	}
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
		node := readNodeFromRedis(context.Background(), key)
		if node != nil { //node is in cache
			log.Printf("Node with key <%s> CACHE HIT", key)
			nodes = append(nodes, graphproto.Node{Key: node.Key, Name: node.Name, Asn: node.Asn, RouterIp: node.Router_ip})
		} else { // node not in cache
			log.Printf("Node with key <%s> CACHE MISS", key)
			var doc NodeDocument
			_, err := col.ReadDocument(ctx, key, &doc)
			if err != nil {
				log.Fatalf("Could not read document with _id: %s, %v", key, err)
			}
			node := graphproto.Node{Key: doc.Key, Name: doc.Name, Asn: doc.Asn, RouterIp: doc.Router_ip}
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func readNodeFromRedis(ctx context.Context, key string) *NodeDocument {
	node := &NodeDocument{}
	bytes, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil //key does not exist
	} else if err != nil {
		panic(err) //error accessing key
	} else {
		err := json.Unmarshal(bytes, node)
		if err != nil {
			log.Fatal("Marshalling error: ", err)
		}
		return node
	}
}

func writeNodeToRedis(key string, node NodeDocument) {
	err := rdb.Set(context.Background(), key, node, 0).Err()
	if err != nil {
		log.Fatalf("Could not write Node to Redis Cache, %v", err)
	}
}

func writeLinkToRedis(key string, link LinkDocument) {
	err := rdb.Set(context.Background(), key, link, 0).Err()
	if err != nil {
		log.Fatalf("Could not write Link to Redis Cache, %v", err)
	}
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
