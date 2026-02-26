package hybrid

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/holodanger/genesis/internal/goservice"
	"github.com/holodanger/genesis/internal/t3"
)

type Config struct {
	ProjectName string
}

func Spawn(rootPath string, projectName string, withAI bool) {
	fmt.Printf("⚔️  [HYBRID] Constructing Twin Architecture: %s | AI: %v\n", projectName, withAI)

	// 1. Create Root Directory (Handled by main, but good to ensure)
	os.MkdirAll(rootPath, 0755)

	// 2. Generate Root Files
	writeRootFiles(rootPath, projectName)

	// 3. Spawn THE SHIELD (Web - T3)
	webPath := filepath.Join(rootPath, "web")
	fmt.Println("  > Spawning Shield Node (Web)...")
	t3.Spawn(webPath, projectName, true) // We name the package the actual project name

	// 3.1 Install Dependencies
	fmt.Println("    📦 [BUN] Installing Dependencies... (Hold Fast)")
	cmd := exec.Command("bun", "install")
	cmd.Dir = webPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("    ⚠️  [WARN] Dependency install failed: %v\n", err)
	}

	// 4. Spawn THE SPEAR (API - Go)
	apiPath := filepath.Join(rootPath, "api")
	fmt.Println("  > Spawning Spear Node (API)...")

	originalWd, _ := os.Getwd()
	os.Chdir(rootPath)

	goBuilder := goservice.NewBuilder("api", withAI)
	if err := goBuilder.Build(); err != nil {
		fmt.Printf("❌ [ERROR] API Spawn failed: %v\n", err)
	}

	os.Chdir(originalWd)

	// 5. THE NEURAL LINK (Rewiring Configs)
	// Both generated projects point to their own DB names (web, api).
	// We must force them to the SHARED TRUTH: {{projectName}}
	fmt.Println("  > Establishing Neural Link (Shared DB Config)...")

	sharedDBUrl := fmt.Sprintf("postgres://postgres:password@localhost:5432/%s", projectName)

	// Rewrite Web .env
	webEnv := fmt.Sprintf(`DATABASE_URL="%s"
BETTER_AUTH_URL="http://localhost:3000"
`, sharedDBUrl)
	os.WriteFile(filepath.Join(webPath, ".env"), []byte(webEnv), 0644)

	// Rewrite API .env
	apiEnv := fmt.Sprintf(`PORT=8080
DATABASE_URL="%s?sslmode=disable"
`, sharedDBUrl)

	if withAI {
		apiEnv += `OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-5-nano
`
	}
	os.WriteFile(filepath.Join(apiPath, ".env"), []byte(apiEnv), 0644)
}

func writeRootFiles(root string, name string) {
	files := map[string]string{
		"compose.yml": RootCompose,
		"Makefile":    RootMakefile,
	}

	config := Config{ProjectName: name}

	for filename, content := range files {
		f, _ := os.Create(filepath.Join(root, filename))
		tmpl, _ := template.New(filename).Parse(content)
		tmpl.Execute(f, config)
		f.Close()
	}
}
