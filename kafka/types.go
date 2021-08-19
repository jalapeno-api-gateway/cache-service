package kafka

type KafkaEventMessage struct {
	TopicType	int 	`json:"TopicType,omitempty"`
	Key			string	`json:"_key,omitempty"`
	Id			string	`json:"_id,omitempty"`
	Action		string	`json:"action,omitempty"`
}