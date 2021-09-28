package kafka

var LSNodeEvents = make(chan KafkaEventMessage)
var LSLinkEvents = make(chan KafkaEventMessage)
var LSPrefixEvents = make(chan KafkaEventMessage)
var LSSRV6SIDEvents = make(chan KafkaEventMessage)