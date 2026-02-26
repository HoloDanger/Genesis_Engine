package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	// IMPORT MODULES
	"github.com/holodanger/genesis/internal/goservice"
	"github.com/holodanger/genesis/internal/hybrid"
	"github.com/holodanger/genesis/internal/t3"
)

func main() {
	// 1. TACTICAL INPUT
	projectName := flag.String("name", "", "Project Name")
	projectType := flag.String("type", "t3", "Archetype: t3 (Frontend Shield) | go (Backend Spear) | hybrid (Twin Architecture)")
	aiEnabled := flag.Bool("ai", false, "Enable AI Features (OpenAI)")
	flag.Parse()

	if *projectName == "" {
		fmt.Println("⚠️  [USAGE] genesis -name <project_name> -type <t3|go|hybrid> -ai=<true|false>")
		os.Exit(1)
	}

	// 2. ROOT ESTABLISHMENT
	// We establish the root path here, but the specific builders
	// handle their internal file structures.
	currentDir, _ := os.Getwd()
	rootPath := filepath.Join(currentDir, *projectName)

	// Create root if it doesn't exist (Idempotency)
	if err := os.MkdirAll(rootPath, 0755); err != nil {
		fmt.Printf("❌ [ERROR] Failed to secure territory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n⚔️  [GENESIS] Spawning Archetype: %s | Node: %s | AI: %v\n", *projectType, *projectName, *aiEnabled)

	// 3. STRATEGY EXECUTION
	switch *projectType {
	case "t3":
		deployT3(rootPath, *projectName)
	case "go":
		deployGoService(*projectName, *aiEnabled)
	case "hybrid":
		hybrid.Spawn(rootPath, *projectName, *aiEnabled)
	default:
		fmt.Printf("❌ [ERROR] Unknown archetype: '%s'. Options: t3, go, hybrid\n", *projectType)
		os.Exit(1)
	}
}

// --- TACTICAL SUBROUTINES ---

func deployT3(root, name string) {
	// 1. Build Files
	t3.Spawn(root, name, false)

	// 2. Install Deps (Bun)
	fmt.Println("📦 [BUN] Installing Dependencies... (Hold Fast)")
	cmd := exec.Command("bun", "install")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  [WARN] Dependency install failed: %v\n", err)
	}

	// 3. Debrief
	printDebrief(name, "bun dev")
}

func deployGoService(name string, withAI bool) {
	// 1. Initialize Builder
	// The Go builder handles its own file generation and 'go mod tidy'
	builder := goservice.NewBuilder(name, withAI)

	// 2. Execute
	if err := builder.Build(); err != nil {
		fmt.Printf("❌ [FATAL] Forge failed: %v\n", err)
		os.Exit(1)
	}

	// 3. Debrief
	printDebrief(name, "make run")
}

func printDebrief(name, runCmd string) {
	fmt.Printf("\n✅ [SUCCESS] Node '%s' is operational.\n", name)
	fmt.Println("   -------------------------------------")
	fmt.Printf("   cd %s\n", name)
	fmt.Printf("   %s\n", runCmd)
	fmt.Println("   -------------------------------------")
}
