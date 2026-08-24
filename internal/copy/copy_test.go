package copy

import "testing"

func TestUserFacingCopyDoesNotContainForbiddenLanguage(t *testing.T) {
	for _, s := range Strings {
		if ContainsForbiddenLanguage(s) {
			t.Errorf("user-facing copy contains forbidden language: %q", s)
		}
	}
}

func TestForbiddenLanguageDetection(t *testing.T) {
	violations := []string{
		"This is a SCAM",
		"report suspected fraud immediately",
		"our same-company policy applies",
		"Same-Company",
	}
	for _, s := range violations {
		if !ContainsForbiddenLanguage(s) {
			t.Errorf("expected forbidden language in %q", s)
		}
	}

	clean := []string{
		AppName,
		AppDescription,
		DisclaimerBody,
	}
	for _, s := range clean {
		if ContainsForbiddenLanguage(s) {
			t.Errorf("unexpected forbidden language in %q", s)
		}
	}
}
