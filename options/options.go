package options

type KafkaOptions struct {
	Broker    string
	Topic     string
	PoolLimit int
}

type Options struct {
	Consumer KafkaOptions
	Producer KafkaOptions
}
