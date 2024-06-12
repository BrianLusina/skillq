package models

import "fmt"

// UserModel represents the model of a user as stored in a database
type UserModel struct {
	BaseModel    BaseModel `bson:"inline"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	Skills       []string  `bson:"skills"`
	ImageUrl     string    `bson:"imageUrl"`
	JobTitle     string    `bson:"jobTitle"`
	PasswordHash string    `bson:"passwordHash"`
}

func (u *UserModel) String() string {
	return fmt.Sprintf("UserModel(base=%s, name=%s, email=%s, skills=%v, imageUrl=%s, jobTitle=%s)",
		u.BaseModel.String(), u.Name, u.Email, u.Skills, u.ImageUrl, u.JobTitle)
}
