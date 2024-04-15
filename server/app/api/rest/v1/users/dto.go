package users

import "time"

// userResponseDto is the DTO for a response on a user request
type userResponseDto struct {
	ID        string     `json:"id,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	XID       string     `json:"xid,omitempty"`
	KeyID     string     `json:"keyId,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Skills    []string   `json:"skills,omitempty"`
	Image     string     `json:"image,omitempty"`
	JobTitle  string     `json:"jobTitle,omitempty"`
}

// userRequestDto is the DTO for a user request
type userRequestDto struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required" validate:"required"`
	Skills   []string `json:"skills" binding:"required"`
	Image    string   `json:"image"`
	JobTitle string   `json:"jobTitle" binding:"required"`
}
