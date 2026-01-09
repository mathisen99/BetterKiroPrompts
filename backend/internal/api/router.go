package api

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", HandleHealth)
	mux.HandleFunc("POST /api/kickoff/generate", HandleKickoffGenerate)
	return mux
}
