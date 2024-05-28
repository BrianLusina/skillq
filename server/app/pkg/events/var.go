package events

type EventName string

const (
	EmailVerificationStartedName EventName = "EmailVerificationStarted"
	EmailVerificationSentName    EventName = "EmailVerificationSent"
)
