package copy

import "testing"

func TestUserFacingCopyDoesNotContainForbiddenLanguage(t *testing.T) {
	for _, s := range Strings {
		if ContainsForbiddenLanguage(s) {
			t.Errorf("user-facing copy contains forbidden language: %q", s)
		}
	}
}

func TestAllowListContainsExactlyRequiredPhrases(t *testing.T) {
	want := []string{
		"reports indicated",
		"appear related",
		"no record found on DATE",
	}
	if len(Strings) != len(want) {
		t.Fatalf("allow-list length = %d, want %d", len(Strings), len(want))
	}
	for i, phrase := range want {
		if Strings[i] != phrase {
			t.Errorf("Strings[%d] = %q, want %q", i, Strings[i], phrase)
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
		ReportsIndicated,
		AppearRelated,
		NoRecordFoundOnDate,
	}
	for _, s := range clean {
		if ContainsForbiddenLanguage(s) {
			t.Errorf("unexpected forbidden language in %q", s)
		}
	}
}
