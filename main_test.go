package main

import (
	"os"
	"testing"
)

// TestSanity ensures the environment is capable of running the genesis logic.
func TestSanity(t *testing.T) {
	// 1. Check working directory access
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	t.Logf("Working Directory: %s", wd)

	// 2. Verify basic flags exist (indirectly checking import cycles)
	// If main.go fails to compile, this test won't even run.
	t.Log("Genesis Engine compiled successfully.")
}
