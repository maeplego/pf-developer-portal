package catalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDir(t *testing.T) {
	dir := filepath.Join("..", "..", "specs")
	c, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.APIs) != 3 {
		t.Fatalf("apis=%d", len(c.APIs))
	}
	pay, ok := c.Get("payments")
	if !ok || pay.Title != "Payments API" {
		t.Fatalf("payments: %+v", pay)
	}
	if len(pay.Paths) < 2 {
		t.Fatalf("payments ops=%d", len(pay.Paths))
	}
}

func TestRejectOversized(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "huge.yaml")
	if err := os.WriteFile(p, make([]byte, maxSpecBytes+1), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(p); err == nil {
		t.Fatal("expected size error")
	}
}
