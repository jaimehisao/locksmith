// Package copy holds the allow-listed user-facing strings for LocksmithWatch SF.
// All user-facing copy in the repository must be registered here.
package copy

// Strings is the authoritative allow-list of user-facing copy.
var Strings = []string{
	AppName,
	AppDescription,
	DisclaimerHeading,
	DisclaimerBody,
	RunningTestsHeading,
	RunningTestsRequirement,
}

const (
	AppName = "LocksmithWatch SF"

	AppDescription = "LocksmithWatch SF is a public repository for monitoring and reporting tools related to locksmith services in San Francisco."

	DisclaimerHeading = "Disclaimer"

	DisclaimerBody = "This project is provided for informational purposes only. It is not legal advice. Consult a qualified attorney for legal questions."

	RunningTestsHeading = "Running tests"

	RunningTestsRequirement = "Requires Go 1.22 or later."
)
