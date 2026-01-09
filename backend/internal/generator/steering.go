package generator

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
	CustomRules        map[string][]string
}

type SteeringFile struct {
	Path    string
	Content string
}

func GenerateSteering(config SteeringConfig) ([]SteeringFile, error) {
	// Implementation in task 24
	return []SteeringFile{}, nil
}
