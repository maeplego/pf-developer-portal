package speclint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSpecsHaveNoSecretLikeExamples(t *testing.T) {
	dir := filepath.Join("..", "..", "specs")
	if err := ScanDir(dir); err != nil {
		t.Fatal(err)
	}
}

func TestDetectsAWSKeyShape(t *testing.T) {
	dir := t.TempDir()
	// Shape only; not a real credential.
	body := "example: AKIA" + "IOSFODNN7EXAMPLE"
	if err := os.WriteFile(filepath.Join(dir, "bad.yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ScanDir(dir); err == nil {
		t.Fatal("expected lint failure")
	}
}
