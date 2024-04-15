package values

import (
	"testing"
)

type emailTestCase struct {
	name      string
	input     string
	shouldErr bool
}

var emailTestCases = []emailTestCase{
	{
		name:      "empty email should return nil and error",
		input:     "",
		shouldErr: true,
	},
	{
		name:      "invalid email should return nil and error",
		input:     "johndoe",
		shouldErr: true,
	},
	{
		name:      "invalid email should return nil and error",
		input:     "johndoe@",
		shouldErr: true,
	},
	{
		name:      "valid email should return email and nil error",
		input:     "johndoe@example.com",
		shouldErr: false,
	},
}

func TestNewEmail(t *testing.T) {
	for _, tc := range emailTestCases {
		t.Run(tc.name, func(t *testing.T) {
			email, err := NewEmail(tc.input)
			if tc.shouldErr {
				if err == nil {
					t.Errorf("NewEmail(%s) = (%v, %v), should return error %v, got %v", tc.input, email, err, tc.shouldErr, err)
				}
			}
		})
	}
}

func BenchmarkNewEmail(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range emailTestCases {
			_, _ = NewEmail(tc.input)
		}
	}
}
