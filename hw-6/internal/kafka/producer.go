package kafka

import "github.com/IBM/sarama"

func NewSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	producerConfig.Producer.RequiredAcks = sarama.WaitForLocal
	producerConfig.Producer.Compression = sarama.CompressionGZIP
	producerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.Return.Errors = true

	return sarama.NewSyncProducer(brokers, producerConfig)
}
