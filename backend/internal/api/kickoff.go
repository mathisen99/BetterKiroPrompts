package api

import (
	"encoding/json"
	"net/http"

	"better-kiro-prompts/internal/generator"
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

	genAnswers := generator.KickoffAnswers{
		ProjectIdentity: req.Answers.ProjectIdentity,
		SuccessCriteria: req.Answers.SuccessCriteria,
		UsersAndRoles:   req.Answers.UsersAndRoles,
		DataSensitivity: req.Answers.DataSensitivity,
		DataLifecycle: generator.DataLifecycle{
			Retention:    req.Answers.DataLifecycle.Retention,
			Deletion:     req.Answers.DataLifecycle.Deletion,
			Export:       req.Answers.DataLifecycle.Export,
			AuditLogging: req.Answers.DataLifecycle.AuditLogging,
			Backups:      req.Answers.DataLifecycle.Backups,
		},
		AuthModel:   req.Answers.AuthModel,
		Concurrency: req.Answers.Concurrency,
		RisksAndTradeoffs: generator.RisksAndTradeoffs{
			TopRisks:    req.Answers.RisksAndTradeoffs.TopRisks,
			Mitigations: req.Answers.RisksAndTradeoffs.Mitigations,
			NotHandled:  req.Answers.RisksAndTradeoffs.NotHandled,
		},
		Boundaries:       req.Answers.Boundaries,
		BoundaryExamples: req.Answers.BoundaryExamples,
		NonGoals:         req.Answers.NonGoals,
		Constraints:      req.Answers.Constraints,
	}

	prompt, err := generator.GenerateKickoff(genAnswers)
	if err != nil {
		http.Error(w, "Failed to generate prompt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(KickoffResponse{Prompt: prompt})
}
