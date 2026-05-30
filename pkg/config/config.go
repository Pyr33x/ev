package config

import (
	"github.com/pyr33x/envy"
)

type Config struct {
	Server    Server
	Kafka     Kafka
	Telemetry Telemetry
}

type Server struct {
	Port int
}

type Conn struct {
	Host string
	Port int
}

type Kafka struct {
	Conn
	GroupID   string
	Topic     string
	Partition int
}

type Telemetry struct {
	ServiceName string
	Environment string
	Endpoint    string
	Insecure    bool
	UseGRPC     bool
}

func New() *Config {
	return &Config{
		Server: Server{
			Port: envy.GetInt("PORT", 80),
		},
		Kafka: Kafka{
			Conn: Conn{
				Host: envy.GetString("KAFKA_HOST", "localhost"),
				Port: envy.GetInt("KAFKA_PORT", 9092),
			},
			GroupID:   envy.GetString("KAFKA_GROUP_ID", "ev_getters"),
			Topic:     envy.GetString("KAFKA_TOPIC", "pipe"),
			Partition: envy.GetInt("KAFKA_PARTITION", 0),
		},
		Telemetry: Telemetry{
			ServiceName: envy.GetString("OTEL_SERVICE_NAME", "ev-service"),
			Environment: envy.GetString("OTEL_ENVIRONMENT", "development"),
			Endpoint:    envy.GetString("OTEL_ENDPOINT", "localhost:4317"),
			Insecure:    envy.GetBool("OTEL_INSECURE", false),
			UseGRPC:     envy.GetBool("OTEL_USEGRPC", false),
		},
	}
}
