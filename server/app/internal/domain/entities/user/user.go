package user

import (
	"fmt"

	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/values"
	"github.com/pkg/errors"
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

	// skillSet is the list of skillSet this user has
	skillSet map[string]bool

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

	// ImageUrl is the URL of the image
	ImageUrl string

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

	skillSet := map[string]bool{}

	skills := params.Skills
	for _, skill := range skills {
		_, ok := skillSet[skill]
		if !ok {
			skillSet[skill] = true
		}
	}

	return User{
		Entity:         entity,
		name:           params.Name,
		email:          *email,
		imageData:      params.ImageData,
		imageUrl:       params.ImageUrl,
		skillSet:       skillSet,
		jobTitle:       params.JobTitle,
		hashedPassword: params.Password,
	}, nil
}

// Name returns the user's name
func (u *User) Name() string {
	return u.name
}

// SetName returns the user's name
func (u *User) SetName(name string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("invalid name %s provided", name)
	}
	u.name = name
	return u, nil
}

// Password retrieves the hashed user password
func (u *User) Password() string {
	return u.hashedPassword
}

// Email returns the user's email address
func (u *User) Email() string {
	return u.email.Get()
}

// SetEmail updates the user's email address
func (u *User) SetEmail(email string) (*User, error) {
	err := u.email.Set(email)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid email %s provided", email)
	}

	return u, nil
}

// ImageUrl returns the user's image URL
func (u *User) ImageData() []byte {
	return u.imageData
}

func (u *User) ImageUrl() string {
	return u.imageUrl
}

func (u *User) SetImageUrl(url string) *User {
	u.imageUrl = url
	return u
}

// WithImage updates the image URL
func (u User) WithImage(imageUrl string) User {
	u.imageUrl = imageUrl
	return u
}

// Skills is a list of all the skills a user has
func (u *User) Skills() []string {
	skills := []string{}
	for skill := range u.skillSet {
		skills = append(skills, skill)
	}
	return skills
}

// SetSkill
func (u *User) SetSkills(skills []string) *User {
	currentSkills := u.skillSet
	for _, skill := range skills {
		_, ok := currentSkills[skill]
		if !ok {
			currentSkills[skill] = true
		}
	}

	u.skillSet = currentSkills

	return u
}

// JobTitle is the user's job title
func (u *User) JobTitle() string {
	return u.jobTitle
}

// SetJobTitle sets the job title and returns a copy of the user
func (u *User) SetJobTitle(title string) (*User, error) {
	if title != "" {
		u.jobTitle = title
		return u, nil
	}
	return u, fmt.Errorf("invalid job title provided: %s", title)
}
