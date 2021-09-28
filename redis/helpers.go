package redis

import "encoding/json"

//Marshaling Methods used to write to Redis

func (node LsNodeDocument) MarshalBinary() ([]byte, error) {
	return json.Marshal(node)
}

func (link LsLinkDocument) MarshalBinary() ([]byte, error) {
	return json.Marshal(link)
}