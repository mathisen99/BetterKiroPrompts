package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"better-kiro-prompts/internal/api"
	"better-kiro-prompts/internal/db"
	"better-kiro-prompts/internal/generation"
	"better-kiro-prompts/internal/openai"
	"better-kiro-prompts/internal/ratelimit"
)

func main() {
	ctx := context.Background()

	if err := db.Connect(ctx); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize dependencies
	var routerCfg *api.RouterConfig

	// Try to create OpenAI client (optional - may not have API key in dev)
	openaiClient, err := openai.NewClient()
	if err != nil {
		log.Printf("Warning: OpenAI client not initialized: %v", err)
		log.Printf("Generation endpoints will not be available")
	} else {
		genService := generation.NewService(openaiClient)
		rateLimiter := ratelimit.NewLimiter()
		routerCfg = &api.RouterConfig{
			GenerationService: genService,
			RateLimiter:       rateLimiter,
		}
		log.Printf("Generation service initialized")
	}

	router := api.NewRouter(routerCfg)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
