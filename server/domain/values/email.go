package values

import "github.com/BrianLusina/skillq/server/utils/validators"

// Email represents a value object that wraps an email address
type Email struct {
	value string
}

// NewEmail creates a new email address after validation, returns an error if email is invalid
func NewEmail(value string) (*Email, error) {
	if err := validators.ValidateEmail(value); err != nil {
		return nil, err
	}

	return &Email{
		value: value,
	}, nil
}

// Get retrieves the email address
func (e *Email) Get() string {
	return e.value
}

// Set, sets the email address
func (e *Email) Set(value string) error {
	err := validators.ValidateEmail(value)
	if err != nil {
		return err
	}
	e.value = value
	return nil
}
