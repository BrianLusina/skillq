package amqp

// Config is the AMQP configuration to establish a connection to an AMQP broker
type Config struct {
	Username string
	Password string
	Host     string
	Port     string
}

// ExchangeOptionParams are the parameters for configuring an exchange
type ExchangeOptionParams struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       map[string]any
}

// QueueOptionParams are the parameters for configuring a queue
type QueueOptionParams struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       map[string]any
}

type ConsumerOptionParams struct {
	Tag       string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      map[string]any
}

type QosOptionParams struct {
	// prefetch count
	PrefetchCount int

	// prefetch size
	PrefetchSize int

	// global
	PrefetchGlobal bool
}
