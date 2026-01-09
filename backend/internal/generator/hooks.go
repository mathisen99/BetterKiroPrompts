package generator

import (
	"encoding/json"
	"strings"

	"better-kiro-prompts/internal/templates"
)

type HooksTechStack struct {
	HasGo         bool
	HasTypeScript bool
	HasReact      bool
}

type HooksConfig struct {
	Preset    string
	TechStack HooksTechStack
}

type HookFile struct {
	Path    string
	Content string
}

type hookDef struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Enabled     bool     `json:"enabled"`
	When        hookWhen `json:"when"`
	Then        hookThen `json:"then"`
}

type hookWhen struct {
	Type string `json:"type"`
}

type hookThen struct {
	Type    string `json:"type"`
	Command string `json:"command,omitempty"`
	Prompt  string `json:"prompt,omitempty"`
}

func GenerateHooks(config HooksConfig) ([]HookFile, error) {
	var files []HookFile

	// Light: formatters
	if config.TechStack.HasGo {
		files = append(files, makeHook("format-go", "Format Go code", "agentStop", "runCommand", "go fmt ./..."))
	}
	if config.TechStack.HasTypeScript || config.TechStack.HasReact {
		files = append(files, makeHook("format-web", "Format frontend code", "agentStop", "runCommand", "pnpm format"))
	}

	if config.Preset == "light" {
		return files, nil
	}

	// Basic: + linters, tests (from templates)
	if config.TechStack.HasGo {
		files = append(files, loadHookTemplate("hooks/lint-go.tmpl", "lint-go"))
	}
	if config.TechStack.HasTypeScript || config.TechStack.HasReact {
		files = append(files, loadHookTemplate("hooks/lint-web.tmpl", "lint-web"))
	}
	files = append(files, loadHookTemplate("hooks/test.tmpl", "run-tests"))

	if config.Preset == "basic" {
		return files, nil
	}

	// Default: + secret scan, prompt guardrails
	files = append(files, makeHookPrompt("secret-scan", "Scan for secrets", "agentStop", "Check for accidentally committed secrets or API keys"))
	files = append(files, makeHookPrompt("prompt-guard", "Prompt guardrails", "promptSubmit", "Review this prompt for unsafe operations and confirm before proceeding"))

	if config.Preset == "default" {
		return files, nil
	}

	// Strict: + static analysis, vuln scan
	if config.TechStack.HasGo {
		files = append(files, makeHook("static-analysis", "Static analysis", "agentStop", "runCommand", "staticcheck ./..."))
	}
	files = append(files, makeHookPrompt("vuln-scan", "Vulnerability scan", "userTriggered", "Run dependency vulnerability scan and report findings"))

	return files, nil
}

func makeHook(name, desc, whenType, thenType, cmd string) HookFile {
	h := hookDef{
		Name: name, Description: desc, Version: "1.0.0", Enabled: true,
		When: hookWhen{Type: whenType},
		Then: hookThen{Type: thenType, Command: cmd},
	}
	content, _ := json.MarshalIndent(h, "", "  ")
	return HookFile{Path: ".kiro/hooks/" + name + ".kiro.hook", Content: string(content)}
}

func makeHookPrompt(name, desc, whenType, prompt string) HookFile {
	h := hookDef{
		Name: name, Description: desc, Version: "1.0.0", Enabled: true,
		When: hookWhen{Type: whenType},
		Then: hookThen{Type: "askAgent", Prompt: prompt},
	}
	content, _ := json.MarshalIndent(h, "", "  ")
	return HookFile{Path: ".kiro/hooks/" + name + ".kiro.hook", Content: string(content)}
}

func loadHookTemplate(path, name string) HookFile {
	content, _ := templates.FS.ReadFile(path)
	return HookFile{Path: ".kiro/hooks/" + name + ".kiro.hook", Content: strings.TrimSpace(string(content))}
}
