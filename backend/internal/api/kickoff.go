package api

import (
	"encoding/json"
	"net/http"
)

type DataLifecycle struct {
	Retention    string `json:"retention"`
	Deletion     string `json:"deletion"`
	Export       string `json:"export"`
	AuditLogging string `json:"auditLogging"`
	Backups      string `json:"backups"`
}

type RisksAndTradeoffs struct {
	TopRisks    []string `json:"topRisks"`
	Mitigations []string `json:"mitigations"`
	NotHandled  []string `json:"notHandled"`
}

type KickoffAnswers struct {
	ProjectIdentity   string            `json:"projectIdentity"`
	SuccessCriteria   string            `json:"successCriteria"`
	UsersAndRoles     string            `json:"usersAndRoles"`
	DataSensitivity   string            `json:"dataSensitivity"`
	DataLifecycle     DataLifecycle     `json:"dataLifecycle"`
	AuthModel         string            `json:"authModel"`
	Concurrency       string            `json:"concurrency"`
	RisksAndTradeoffs RisksAndTradeoffs `json:"risksAndTradeoffs"`
	Boundaries        string            `json:"boundaries"`
	BoundaryExamples  []string          `json:"boundaryExamples"`
	NonGoals          string            `json:"nonGoals"`
	Constraints       string            `json:"constraints"`
}

type KickoffRequest struct {
	Answers KickoffAnswers `json:"answers"`
}

type KickoffResponse struct {
	Prompt string `json:"prompt"`
}

func HandleKickoffGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req KickoffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Answers.ProjectIdentity == "" {
		http.Error(w, "projectIdentity is required", http.StatusBadRequest)
		return
	}

	// TODO: Call generator.GenerateKickoff(req.Answers)
	prompt := "Generated kickoff prompt placeholder"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(KickoffResponse{Prompt: prompt})
}
