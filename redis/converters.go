package redis

import "github.com/jalapeno-api-gateway/cache-service/arangodb"

func ConvertToRedisLsNode(document arangodb.LsNodeDocument) LsNodeDocument {
	return LsNodeDocument{
		Id: document.Id,
		Key: document.Key,
		Name: document.Name,
		Asn: document.Asn,
		Router_ip: document.Router_ip,
	}
}

func ConvertToRedisLsLink(document arangodb.LsLinkDocument) LsLinkDocument {
	return LsLinkDocument{
		Id: document.Id,
		Key: document.Key,
		Router_ip: document.Router_ip,
		Peer_ip: document.Peer_ip,
		LocalLink_ip: document.LocalLink_ip,
		RemoteLink_ip: document.RemoteLink_ip,
		Igp_metric: document.Igp_metric,
	}
}