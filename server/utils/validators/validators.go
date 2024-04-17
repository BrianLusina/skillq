package validators

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	ozzoValidation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-playground/validator/v10"
)

var (
	// TODO: update pattern to accommodate more phone numbers
	// Consider github.com/ttacon/libphonenumber
	phoneNumberRegexPattern = `(((\+)*254)|0)7[0-9]{8}`
	phoneNumberRegex        = regexp.MustCompile(phoneNumberRegexPattern)
)

// ValidateUrl checks if the provided URL is a valid URL
func ValidateUrl(url string) error {
	err := ozzoValidation.Validate(url, ozzoValidation.Required, is.URL)
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
	err := ozzoValidation.Validate(email, ozzoValidation.Required, is.Email)
	if err != nil {
		return fmt.Errorf("provided email address: %s is invalid", email)
	}
	return nil
}

// ValidateCountryCode validates a country code
func ValidateCountryCode(country string) error {
	err := ozzoValidation.Validate(country, is.CountryCode3)

	if err != nil {
		return fmt.Errorf("provided country code %s is not a valid ISO3166 Alpha 3 country code", country)
	}

	return nil
}

// ValidateCurrencyCode validates that a given currency code follows ISO 4217 code
func ValidateCurrencyCode(code string) error {
	err := ozzoValidation.Validate(code, is.CurrencyCode)

	if err != nil {
		return fmt.Errorf("provided currency code %s is not a valid ISO4217 currency code", code)
	}

	return nil
}

// GetValidationErrMsg checks to see if the provided err is a validation error and
// returns the first validation error message.
func GetValidationErrMsg(s any, err error) (errMsg string) {
	fieldErrors := validator.ValidationErrors{}

	if ok := errors.As(err, &fieldErrors); ok {
		fieldErr := fieldErrors[0]
		fieldName := getStructTag(s, fieldErr.Field(), "json")

		switch fieldErr.Tag() {
		case "required":
			errMsg = fmt.Sprintf("%s is a required field", fieldName)
		default:
			errMsg = fmt.Sprintf("Invalid input on %s", fieldName)
		}
	}

	return errMsg
}

func getStructTag(s any, fieldName string, tagKey string) string {
	t := reflect.TypeOf(s)
	field, found := t.FieldByName(fieldName)

	if t.Kind() != reflect.Struct {
		return fieldName
	}

	if !found {
		return fieldName
	}

	return field.Tag.Get(tagKey)
}

// IsValidationError checks to see if error is of type validator.ValidationErrors.
func IsValidationError(err error) bool {
	return errors.As(err, &validator.ValidationErrors{})
}
