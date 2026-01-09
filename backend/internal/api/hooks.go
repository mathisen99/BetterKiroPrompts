package api

import (
	"encoding/json"
	"net/http"

	"better-kiro-prompts/internal/generator"
)

type HooksTechStack struct {
	HasGo         bool `json:"hasGo"`
	HasTypeScript bool `json:"hasTypeScript"`
	HasReact      bool `json:"hasReact"`
}

type HooksRequest struct {
	Preset    string         `json:"preset"`
	TechStack HooksTechStack `json:"techStack"`
}

type HooksResponse struct {
	Files []GeneratedFile `json:"files"`
}

func HandleHooksGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req HooksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Preset == "" {
		http.Error(w, "preset is required", http.StatusBadRequest)
		return
	}

	genConfig := generator.HooksConfig{
		Preset: req.Preset,
		TechStack: generator.HooksTechStack{
			HasGo:         req.TechStack.HasGo,
			HasTypeScript: req.TechStack.HasTypeScript,
			HasReact:      req.TechStack.HasReact,
		},
	}

	files, err := generator.GenerateHooks(genConfig)
	if err != nil {
		http.Error(w, "Failed to generate hooks", http.StatusInternalServerError)
		return
	}

	resp := HooksResponse{Files: make([]GeneratedFile, len(files))}
	for i, f := range files {
		resp.Files[i] = GeneratedFile{Path: f.Path, Content: f.Content}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
