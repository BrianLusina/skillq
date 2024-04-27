package models

// UserModel represents the model of a user as stored in a database
type UserModel struct {
	BaseModel    BaseModel `bson:"inline"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	Skills       []string  `bson:"skills"`
	ImageUrl     string    `bson:"image_url"`
	JobTitle     string    `bson:"job_title"`
	PasswordHash string    `bson:"password_hash"`
}
