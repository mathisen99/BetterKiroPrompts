package templates

import "embed"

//go:embed *.tmpl steering/*.tmpl hooks/*.tmpl
var FS embed.FS
