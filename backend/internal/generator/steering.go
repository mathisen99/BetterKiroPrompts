package generator

import (
	"bytes"
	"text/template"

	"better-kiro-prompts/internal/templates"
)

type TechStack struct {
	Backend  string
	Frontend string
	Database string
}

type SteeringConfig struct {
	ProjectName        string
	ProjectDescription string
	TechStack          TechStack
	IncludeConditional bool
	IncludeManual      bool
	CustomRules        map[string][]string
}

type SteeringFile struct {
	Path    string
	Content string
}

var foundationTemplates = []struct {
	tmpl string
	path string
}{
	{"steering/product.tmpl", ".kiro/steering/product.md"},
	{"steering/tech.tmpl", ".kiro/steering/tech.md"},
	{"steering/structure.tmpl", ".kiro/steering/structure.md"},
}

var conditionalTemplates = []struct {
	tmpl string
	path string
}{
	{"steering/security-go.tmpl", ".kiro/steering/security-go.md"},
	{"steering/security-web.tmpl", ".kiro/steering/security-web.md"},
	{"steering/quality-go.tmpl", ".kiro/steering/quality-go.md"},
	{"steering/quality-web.tmpl", ".kiro/steering/quality-web.md"},
}

func GenerateSteering(config SteeringConfig) ([]SteeringFile, error) {
	var files []SteeringFile

	// Foundation files
	for _, t := range foundationTemplates {
		content, err := renderTemplate(t.tmpl, config)
		if err != nil {
			return nil, err
		}
		files = append(files, SteeringFile{Path: t.path, Content: content})
	}

	// AGENTS.md
	agentsContent, err := renderTemplate("steering/agents.tmpl", config)
	if err != nil {
		return nil, err
	}
	files = append(files, SteeringFile{Path: "AGENTS.md", Content: agentsContent})

	// Conditional files
	if config.IncludeConditional {
		for _, t := range conditionalTemplates {
			content, err := renderTemplate(t.tmpl, config)
			if err != nil {
				return nil, err
			}
			files = append(files, SteeringFile{Path: t.path, Content: content})
		}
	}

	// Manual files
	if config.IncludeManual {
		content, err := renderTemplate("steering/manual-example.tmpl", config)
		if err != nil {
			return nil, err
		}
		files = append(files, SteeringFile{Path: ".kiro/steering/manual-example.md", Content: content})
	}

	return files, nil
}

func renderTemplate(name string, data any) (string, error) {
	tmplContent, err := templates.FS.ReadFile(name)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(name).Parse(string(tmplContent))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
