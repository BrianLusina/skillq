package amqp

// Config is the AMQP configuration to establish a connection to an AMQP broker
type Config struct {
	Username string
	Password string
	Host     string
	Port     string
}
