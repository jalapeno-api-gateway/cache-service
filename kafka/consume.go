package kafka

import (
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

func newSaramaConsumer() sarama.Consumer {
	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("http://%s", os.Getenv("KAFKA_ADDRESS"))}, sarama.NewConfig())
	if err != nil {
		panic(err)
	}
	return consumer
}

func newPartitionConsumer(consumer sarama.Consumer, topic string) sarama.PartitionConsumer {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	return partitionConsumer
}

func closeConsumers(consumer sarama.Consumer, partitionConsumers ...sarama.PartitionConsumer) {
	if err := consumer.Close(); err != nil {
		log.Fatalln(err)
	}

	for _, c := range partitionConsumers {
		if err := c.Close(); err != nil {
			log.Fatalln(err)
		}
	}
}