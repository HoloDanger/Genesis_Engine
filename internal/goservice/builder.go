package goservice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type Builder struct {
	Name   string
	WithAI bool
}

func NewBuilder(name string, withAI bool) *Builder {
	return &Builder{
		Name:   name,
		WithAI: withAI,
	}
}

func (b *Builder) Build() error {
	// 1. Define the File Map
	// Mapped to the constants in templates.go
	files := map[string]string{
		"go.mod":                        GoMod,
		"cmd/api/main.go":               MainGo,
		"internal/config/config.go":     ConfigGo,
		"internal/server/server.go":     ServerGo,
		"internal/server/routes.go":     RoutesGo,
		"internal/server/middleware.go": MiddlewareGo,
		"internal/server/audit.go":      AuditGo,
		"internal/db/schema.sql":        SchemaSQL,
		"internal/db/query.sql":         QuerySQL,
		"sqlc.yaml":                     SQLCConfig,
		".golangci.yml":                 LintConfig,
		"compose.yml":                   DockerCompose,
		"Makefile":                      Makefile,
		".env":                          EnvFile,
		".gitignore":                    GitIgnore,
	}

	// Logic Gate: Inject AI Modules
	if b.WithAI {
		files["internal/ai/service.go"] = AIServiceGo
		files["internal/ai/handler.go"] = AIHandlerGo
	}

	// 2. Data for Templates
	data := map[string]interface{}{
		"Name":   b.Name,
		"WithAI": b.WithAI,
	}

	// 3. Execution Loop
	for path, content := range files {
		fullPath := filepath.Join(b.Name, path)

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("mkdir failed: %w", err)
		}

		// Create File
		f, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("create file failed: %w", err)
		}

		// Parse & Execute
		tmpl, err := template.New(path).Parse(content)
		if err != nil {
			f.Close()
			return fmt.Errorf("parse template failed: %w", err)
		}

		if err := tmpl.Execute(f, data); err != nil {
			f.Close()
			return fmt.Errorf("execute template failed: %w", err)
		}
		f.Close()
	}

	// 4. Formatting & Initialization
	fmt.Println("⚡ [SYSTEM] Initializing Go Module...")

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = b.Name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Warning: 'go mod tidy' failed: %v\n", err)
	}

	return nil
}
