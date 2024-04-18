package user

import (
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/values"
)

// User structure represents a user entity in the system
type User struct {
	entity.Entity

	// name is the user's name
	name string

	// email of the user
	email values.Email

	// hashedPassword is the user's hashed password
	hashedPassword string

	// imageData contains the image Bytes which is the actual image
	imageData []byte

	// imageUrl is the URL to the image
	imageUrl string

	// skills is the list of skills this user has
	skills []string

	// jobTitle is the user's job title
	jobTitle string
}

type UserParams struct {
	// EntityParams contain common parameters for an entity
	entity.EntityParams

	// Name is the user's name
	Name string

	// Email is the user's email address
	Email string

	// Password is the hashed password when creating a user
	Password string

	// ImageData is the image URL for the user
	ImageData []byte

	// Skills is the list of skills this user has
	Skills []string

	// JobTitle is a user's job title
	JobTitle string
}

// New creates a new user entity & potentially an error
func New(params UserParams) (User, error) {
	entity := entity.NewEntity(params.EntityParams)
	email, err := values.NewEmail(params.Email)
	if err != nil {
		return User{}, err
	}

	return User{
		Entity:         entity,
		name:           params.Name,
		email:          *email,
		imageData:      params.ImageData,
		skills:         params.Skills,
		jobTitle:       params.JobTitle,
		hashedPassword: params.Password,
	}, nil
}

// Name returns the user's name
func (u *User) Name() string {
	return u.name
}

// Email returns the user's email address
func (u *User) Email() string {
	return u.email.Get()
}

// ImageUrl returns the user's image URL
func (u *User) ImageData() []byte {
	return u.imageData
}

func (u *User) ImageUrl() string {
	return u.imageUrl
}

// Skills is a list of all the skills a user has
func (u *User) Skills() []string {
	return u.skills
}

// JobTitle is the user's job title
func (u *User) JobTitle() string {
	return u.jobTitle
}
