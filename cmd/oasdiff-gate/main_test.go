package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGateExitCodes(t *testing.T) {
	root := filepath.Join("..", "..")
	run := func(base, rev string) *exec.Cmd {
		cmd := exec.Command("go", "run", ".", filepath.Join(root, "testdata", "openapi", base), filepath.Join(root, "testdata", "openapi", rev))
		cmd.Dir = "."
		cmd.Env = append(os.Environ(), "GOWORK=off")
		return cmd
	}
	if err := run("base.yaml", "compatible.yaml").Run(); err != nil {
		t.Fatalf("compatible should exit 0: %v", err)
	}
	err := run("base.yaml", "breaking.yaml").Run()
	if err == nil {
		t.Fatal("breaking should exit 1")
	}
	if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 {
		t.Fatalf("want exit 1 got %v", err)
	}
}
