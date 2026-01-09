package api

import (
	"encoding/json"
	"net/http"

	"better-kiro-prompts/internal/generator"
)

type TechStack struct {
	Backend  string `json:"backend"`
	Frontend string `json:"frontend"`
	Database string `json:"database"`
}

type SteeringConfig struct {
	ProjectName        string              `json:"projectName"`
	ProjectDescription string              `json:"projectDescription"`
	TechStack          TechStack           `json:"techStack"`
	IncludeConditional bool                `json:"includeConditional"`
	IncludeManual      bool                `json:"includeManual"`
	FileReferences     []string            `json:"fileReferences"`
	CustomRules        map[string][]string `json:"customRules"`
}

type SteeringRequest struct {
	Config SteeringConfig `json:"config"`
}

type GeneratedFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type SteeringResponse struct {
	Files []GeneratedFile `json:"files"`
}

func HandleSteeringGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SteeringRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Config.ProjectName == "" {
		http.Error(w, "projectName is required", http.StatusBadRequest)
		return
	}

	genConfig := generator.SteeringConfig{
		ProjectName:        req.Config.ProjectName,
		ProjectDescription: req.Config.ProjectDescription,
		TechStack: generator.TechStack{
			Backend:  req.Config.TechStack.Backend,
			Frontend: req.Config.TechStack.Frontend,
			Database: req.Config.TechStack.Database,
		},
		IncludeConditional: req.Config.IncludeConditional,
		IncludeManual:      req.Config.IncludeManual,
		FileReferences:     req.Config.FileReferences,
		CustomRules:        req.Config.CustomRules,
	}

	files, err := generator.GenerateSteering(genConfig)
	if err != nil {
		http.Error(w, "Failed to generate steering files", http.StatusInternalServerError)
		return
	}

	resp := SteeringResponse{Files: make([]GeneratedFile, len(files))}
	for i, f := range files {
		resp.Files[i] = GeneratedFile{Path: f.Path, Content: f.Content}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
