package options

type KafkaOptions struct {
	Broker    string
	Topic     string
	PoolLimit int
}

type Options struct {
	Consumer KafkaOptions // Consumer provides configuration for kafka consumer.
	Producer KafkaOptions // Producer provides configuration for kafka producer.
}
