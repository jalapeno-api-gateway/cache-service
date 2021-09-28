package kafka

import (
	"os"
)

func StartEventConsumption() {
	consumer := newSaramaConsumer()
	lsNodeEventsConsumer := newPartitionConsumer(consumer, os.Getenv("LSNODE_KAFKA_TOPIC"))
	lsLinkEventsConsumer := newPartitionConsumer(consumer, os.Getenv("LSLINK_KAFKA_TOPIC"))
	lsPrefixEventsConsumer := newPartitionConsumer(consumer, os.Getenv("LSPREFIX_KAFKA_TOPIC"))
	lsSRV6SIDEventsConsumer := newPartitionConsumer(consumer, os.Getenv("LSSRV6SID_KAFKA_TOPIC"))

	go func() {	
		defer func() {
			closeConsumers(consumer, lsNodeEventsConsumer, lsLinkEventsConsumer)
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