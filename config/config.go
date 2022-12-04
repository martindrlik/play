package config

type Config struct {
	KafkaBroker      string
	KafkaUploadTopic string

	// RequestLimit is limit of in-flight requests that server can handle.
	RequestLimit uint

	ApiKeys []ApiKey
}

type ApiKey struct {
	Name  string
	Value string
}
