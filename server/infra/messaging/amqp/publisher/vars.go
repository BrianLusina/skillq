package publisher

const (
	_retryTimes     = 5
	_backOffSeconds = 2

	_publishMandatory = false
	_publishImmediate = false
	_exchangeName     = "skillq-exchange"
	_bindingKey       = "skillq-routing-key"

	_messageTypeName = "skillq"
)
