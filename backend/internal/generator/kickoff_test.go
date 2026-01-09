package generator

import (
	"strings"
	"testing"
)

func TestGenerateKickoff_ValidAnswers(t *testing.T) {
	answers := KickoffAnswers{
		ProjectIdentity: "Task management app",
		SuccessCriteria: "Users can create and complete tasks",
		UsersAndRoles:   "Anonymous, authenticated users, admins",
		DataSensitivity: "User emails (PII), hashed passwords",
		DataLifecycle: DataLifecycle{
			Retention: "2 years",
			Deletion:  "Soft delete",
		},
		AuthModel:   "basic",
		Concurrency: "Multi-user access",
		RisksAndTradeoffs: RisksAndTradeoffs{
			TopRisks:    []string{"Data loss"},
			Mitigations: []string{"Daily backups"},
		},
		Boundaries:       "Task details are private",
		BoundaryExamples: []string{"Users can only see their own tasks"},
		NonGoals:         "Mobile app",
		Constraints:      "Ship in 2 weeks",
	}

	result, err := GenerateKickoff(answers)
	if err != nil {
		t.Fatalf("GenerateKickoff failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Verify key fields appear in output
	checks := []string{
		"Task management app",
		"Users can create and complete tasks",
		"basic",
		"Mobile app",
	}
	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("Expected output to contain %q", check)
		}
	}
}

func TestGenerateKickoff_EmptyAnswers(t *testing.T) {
	answers := KickoffAnswers{}

	result, err := GenerateKickoff(answers)
	if err != nil {
		t.Fatalf("GenerateKickoff failed with empty answers: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result even with empty answers")
	}
}
