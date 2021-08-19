package arangodb

type LsNodeDocument struct {
	Id			string	`json:"_id,omitempty"`
	Key			string	`json:"_key,omitempty"`
	Name		string	`json:"name,omitempty"`
	Asn			int32	`json:"asn,omitempty"`
	Router_ip	string	`json:"router_ip,omitempty"`
}

type LsLinkDocument struct {
	Id				string	`json:"_id,omitempty"`
	Key				string	`json:"_key,omitempty"`
	Router_ip		string	`json:"router_ip,omitempty"`
	Peer_ip			string	`json:"peer_ip,omitempty"`
	LocalLink_ip	string	`json:"local_link_ip,omitempty"`
	RemoteLink_ip	string	`json:"remote_link_ip,omitempty"`
	Igp_metric		int		`json:"igp_metric,omitempty"`
}