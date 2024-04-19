package user

import (
	"time"

	"github.com/BrianLusina/skillq/server/domain/id"
)

// UserVerification is a structure that contains user verification details
type UserVerification struct {
	id         id.UUID
	code       string
	userId     id.UUID
	isVerified bool
	createdAt  time.Time
	updatedAt  time.Time
}

// UserVerificationParams defines a structure with fields used to create a user verification struct
type UserVerificationParams struct {
	ID         id.UUID
	Code       string
	UserId     id.UUID
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewVerification creates a new user verification structure from a given params
func NewVerification(params UserVerificationParams) UserVerification {
	return UserVerification{
		id:         params.ID,
		code:       params.Code,
		userId:     params.UserId,
		isVerified: params.IsVerified,
		createdAt:  params.CreatedAt,
		updatedAt:  params.UpdatedAt,
	}
}

// ID retrieves the ID of the verification
func (v *UserVerification) ID() id.UUID {
	return v.id
}

// Code retrieves the code of the verification
func (v *UserVerification) Code() string {
	return v.code
}

// UserID retrieves the user ID of the verification
func (v *UserVerification) UserID() id.UUID {
	return v.userId
}

// IsVerified retrieves the verification status
func (v *UserVerification) IsVerified() bool {
	return v.isVerified
}

// CreatedAt retrieves the created at timestamp of the verification
func (v *UserVerification) CreatedAt() time.Time {
	return v.createdAt
}

// UpdatedAt retrieves the updated at timestamp of the verification
func (v *UserVerification) UpdatedAt() time.Time {
	return v.updatedAt
}
