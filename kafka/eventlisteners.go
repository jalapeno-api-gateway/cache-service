package kafka

func StartEventConsumption() {
	consumer := newSaramaConsumer()
	lsNodeEventsConsumer := newPartitionConsumer(consumer, LSNODE_KAFKA_TOPIC)
	lsLinkEventsConsumer := newPartitionConsumer(consumer, LSLINK_KAFKA_TOPIC)
	lsPrefixEventsConsumer := newPartitionConsumer(consumer, LSPREFIX_KAFKA_TOPIC)
	lsSRV6SIDEventsConsumer := newPartitionConsumer(consumer, LSSRV6SID_KAFKA_TOPIC)

	go func() {	
		defer func() {
			closeConsumers(
				consumer,
				lsNodeEventsConsumer,
				lsLinkEventsConsumer,
				lsPrefixEventsConsumer,
				lsSRV6SIDEventsConsumer,
			)
		}()
		
		for {
			select {
			case msg := <-lsNodeEventsConsumer.Messages():
				LSNodeEvents <- unmarshalKafkaMessage(msg)
			case msg := <-lsLinkEventsConsumer.Messages():
				LSLinkEvents <- unmarshalKafkaMessage(msg)
			case msg := <-lsPrefixEventsConsumer.Messages():
				LSPrefixEvents <- unmarshalKafkaMessage(msg)
			case msg := <-lsSRV6SIDEventsConsumer.Messages():
				LSSRV6SIDEvents <- unmarshalKafkaMessage(msg)
			}
		}
	}()
}