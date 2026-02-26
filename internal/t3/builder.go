package t3

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Config struct {
	Name     string
	IsHybrid bool
}

func Spawn(rootPath string, name string, isHybrid bool) {
	fmt.Println("  [T3] Injecting Next.js 16 Architecture...")

	// 1. Create Directories
	dirs := []string{
		"src/app",
		"src/server",
		"src/lib",
		"src/app/api/auth/[...all]",
		"src/app/auth",
	}

	for _, dir := range dirs {
		path := filepath.Join(rootPath, dir)
		os.MkdirAll(path, 0755)
	}

	// 2. Write the Files
	// We use a helper map to link Filenames -> Template Content
	files := map[string]string{
		"package.json":                       PackageJSON,
		"tsconfig.json":                      TSConfig,
		"tailwind.config.ts":                 TailwindConfig,
		"postcss.config.mjs":                 PostCSSConfig,
		"drizzle.config.ts":                  DrizzleConfig,
		".env":                               EnvFile,
		"next.config.ts":                     NextConfig,
		".gitignore":                         GitIgnore,
		"src/app/globals.css":                GlobalCSS,
		"src/app/layout.tsx":                 RootLayout,
		"src/app/page.tsx":                   MainPage,
		"src/server/db.ts":                   DBConnection,
		"src/server/schema.ts":               DBSchema,
		"src/lib/auth.ts":                    AuthServer,
		"src/lib/auth-client.ts":             AuthClient,
		"src/app/api/auth/[...all]/route.ts": AuthRoute,
		"src/app/auth/page.tsx":              LoginPage,
		"compose.yml":                        DockerCompose,
	}

	config := Config{
		Name:     name,
		IsHybrid: isHybrid,
	}

	for filename, content := range files {
		fullPath := filepath.Join(rootPath, filename)

		// Create the file
		f, _ := os.Create(fullPath)

		// Parse and Execute Template (in case we need {{.Name}})
		tmpl, _ := template.New(filename).Parse(content)
		tmpl.Execute(f, config)
		f.Close()

		fmt.Printf("    ├── Injected: %s\n", filename)
	}
}
