package tasks

type TaskName string

const (
	StoreUserImageTaskName     TaskName = "StoreUserImage"
	SendEmailVerificationName  TaskName = "SendEmailVerification"
	StartEmailVerificationName TaskName = "StartEmailVerification"
)
