package v1

import "time"

// userResponseDto is the DTO for a response on a user request
type userResponseDto struct {
	UUID      string     `json:"uuid"`
	XID       string     `json:"xid"`
	KeyID     string     `json:"keyId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	JobTitle  string     `json:"jobTitle"`
	Skills    []string   `json:"skills"`
	ImageUrl  string     `json:"imageUrl"`
}

// userRequestDto is the DTO for a user request
type userRequestDto struct {
	Name     string       `json:"name" binding:"required" validate:"required,min=2,max=24"`
	Email    string       `json:"email" binding:"required" validate:"email,required"`
	Skills   []string     `json:"skills" binding:"required"`
	Image    userImageDto `json:"image"`
	JobTitle string       `json:"jobTitle" binding:"required"`
}

// userImageDto is the DTO for a user image
type userImageDto struct {
	ImageType string `json:"type" validate:"required"`
	Content   string `json:"content" validate:"required"`
}
