package tasks

type TaskName string

const (
	StoreUserImageTaskName     TaskName = "StoreUserImageTask"
	SendEmailVerificationName  TaskName = "SendEmailVerificationTask"
	StartEmailVerificationName TaskName = "StartEmailVerificationTask"
)
