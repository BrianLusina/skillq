package templates

import "os"

func BuildEmailVerification(to string, name string, code string) []byte {
	frontendURL := os.Getenv("FRONTEND_URL")
	link := frontendURL + "/verify-email?code=" + code
	msg := []byte("To: " + to + "\r\n" +
		"Subject: SkillQ: Verify your email address\r\n" +
		"\r\n" +
		"Hi " + name + ",\n\n" +
		"Please follow the link to verify your account: " + link + "\r\n")

	return msg
}
