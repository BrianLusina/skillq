package models

// UserModel represents
type UserModel struct {
	BaseModel BaseModel `bson:"inline"`
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Skills    []string  `json:"skills,omitempty" bson:"skills,omitempty"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty"`
	JobTitle  string    `json:"jobTitle,omitempty" bson:"jobTitle,omitempty"`
}
