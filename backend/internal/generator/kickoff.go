package generator

import (
	"bytes"
	"text/template"

	"better-kiro-prompts/internal/templates"
)

type DataLifecycle struct {
	Retention    string
	Deletion     string
	Export       string
	AuditLogging string
	Backups      string
}

type RisksAndTradeoffs struct {
	TopRisks    []string
	Mitigations []string
	NotHandled  []string
}

type KickoffAnswers struct {
	ProjectIdentity   string
	SuccessCriteria   string
	UsersAndRoles     string
	DataSensitivity   string
	DataLifecycle     DataLifecycle
	AuthModel         string
	Concurrency       string
	RisksAndTradeoffs RisksAndTradeoffs
	Boundaries        string
	BoundaryExamples  []string
	NonGoals          string
	Constraints       string
}

func GenerateKickoff(answers KickoffAnswers) (string, error) {
	tmplContent, err := templates.FS.ReadFile("kickoff.tmpl")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("kickoff").Parse(string(tmplContent))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, answers); err != nil {
		return "", err
	}

	return buf.String(), nil
}
