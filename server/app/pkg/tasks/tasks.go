package tasks

// StoreUserImage is a task message that triggers the storage of a user image
type StoreUserImage struct {
	UserUUID    string `json:"userUUID"`
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
	Name        string `json:"name"`
	Bucket      string `json:"bucket"`
}
