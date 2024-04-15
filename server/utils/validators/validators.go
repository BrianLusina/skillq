// Package validators contain helper functions for validation of fields and values
package validators

import (
	"fmt"
	"regexp"

	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var (
	// TODO: update pattern to accommodate more phone numbers
	// Consider github.com/ttacon/libphonenumber
	phoneNumberRegexPattern = `(((\+)*254)|0)7[0-9]{8}`
	phoneNumberRegex        = regexp.MustCompile(phoneNumberRegexPattern)
)

// ValidateUrl checks if the provided URL is a valid URL
func ValidateUrl(url string) error {
	err := validator.Validate(url, validator.Required, is.URL)
	if err != nil {
		return fmt.Errorf("provided url: %s is invalid", url)
	}
	return nil
}

// ValidatePhoneNumber validates the provided phone number.
func ValidatePhoneNumber(phoneNumber string) error {
	if !phoneNumberRegex.MatchString(phoneNumber) {
		return fmt.Errorf("provided phone number: %s is invalid", phoneNumber)
	}
	return nil
}

// ValidateEmail validates an email address returning an error if invalid
func ValidateEmail(email string) error {
	err := validator.Validate(email, validator.Required, is.Email)
	if err != nil {
		return fmt.Errorf("provided email address: %s is invalid", email)
	}
	return nil
}

// ValidateCountryCode validates a country code
func ValidateCountryCode(country string) error {
	err := validator.Validate(country, is.CountryCode3)

	if err != nil {
		return fmt.Errorf("provided country code %s is not a valid ISO3166 Alpha 3 country code", country)
	}

	return nil
}

// ValidateCurrencyCode validates that a given currency code follows ISO 4217 code
func ValidateCurrencyCode(code string) error {
	err := validator.Validate(code, is.CurrencyCode)

	if err != nil {
		return fmt.Errorf("provided currency code %s is not a valid ISO4217 currency code", code)
	}

	return nil
}
