package amqppublisher

const (
	_publishMandatory = false
	_publishImmediate = false
	_exchangeName     = "skillq-exchange"
	_exchangeKind     = "fanout"
	_bindingKey       = "skillq-routing-key"

	_messageTypeName = "skillq"
)
