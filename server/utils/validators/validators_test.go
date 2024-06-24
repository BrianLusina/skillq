package validators

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	input       string
	expectedErr error
}

var urlTestCases = []testCase{
	{
		name:        "valid url with protocol",
		input:       "https://www.google.com",
		expectedErr: nil,
	},
	{
		name:        "valid url without protocol",
		input:       "www.google.com",
		expectedErr: nil,
	},
	{
		name:        "invalid url with invalid protocol",
		input:       "htt://www.google.com",
		expectedErr: errors.New("provided url: htt://www.google.com is invalid"),
	},
}

func TestValidateUrl(t *testing.T) {
	t.Parallel()

	for _, tc := range urlTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUrl(tc.input)
			if tc.expectedErr != nil {
				assert.Error(t, err)
			}
		})
	}
}

func BenchmarkValidateUrl(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range urlTestCases {
			_ = ValidateUrl(tc.input)
		}
	}
}

var phoneNumberTestCases = []testCase{
	{
		name:        "valid phone number with country code +254700000000",
		input:       "+254700000000",
		expectedErr: nil,
	},
	{
		name:        "valid phone number without country code 0700000000",
		input:       "0700000000",
		expectedErr: nil,
	},
	{
		name:        "invalid phone number without country code or leading 0 700000000",
		input:       "700000000",
		expectedErr: errors.New("provided phone number: 700000000 is invalid"),
	},
}

func TestValidatePhoneNumber(t *testing.T) {
	t.Parallel()

	for _, tc := range phoneNumberTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePhoneNumber(tc.input)
			if tc.expectedErr != nil {
				assert.Error(t, err)
			}
		})
	}
}

func BenchmarkValidatePhoneNumber(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range phoneNumberTestCases {
			_ = ValidatePhoneNumber(tc.input)
		}
	}
}

var countryCodeTestCases = []testCase{
	{
		name:        "valid country code KE",
		input:       "KE",
		expectedErr: nil,
	},
	{
		name:        "valid country code ke",
		input:       "ke",
		expectedErr: nil,
	},
	{
		name:        "invalid country code K3N",
		input:       "K3N",
		expectedErr: errors.New("provided country code is invalid"),
	},
}

func TestValidateCountryCode(t *testing.T) {
	t.Parallel()

	for _, tc := range countryCodeTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCountryCode(tc.input)
			if tc.expectedErr != nil {
				assert.Error(t, err)
			}
		})
	}
}

func BenchmarkValidateCountryCode(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range countryCodeTestCases {
			_ = ValidateCountryCode(tc.input)
		}
	}
}
