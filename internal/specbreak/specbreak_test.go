package specbreak

import (
	"os"
	"path/filepath"
	"testing"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "testdata", "openapi", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestCompatibleOptionalFields(t *testing.T) {
	rep, err := CompareYAML(read(t, "base.yaml"), read(t, "compatible.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if rep.HasERR() {
		t.Fatalf("unexpected errors %+v", rep.Errors)
	}
}

func TestBreakingFieldAndPathRemoval(t *testing.T) {
	rep, err := CompareYAML(read(t, "base.yaml"), read(t, "breaking.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !rep.HasERR() {
		t.Fatal("expected ERR")
	}
	want := map[string]bool{
		"api-path-removed":                    false,
		"response-property-removed":           false,
		"request-property-became-required":    false,
	}
	for _, e := range rep.Errors {
		if _, ok := want[e.ID]; ok {
			want[e.ID] = true
		}
	}
	for id, ok := range want {
		if !ok {
			t.Fatalf("missing %s in %+v", id, rep.Errors)
		}
	}
	var sawAmount bool
	for _, e := range rep.Errors {
		if e.ID == "response-property-removed" && e.Message == "response 201 property removed: amountMinor" {
			sawAmount = true
		}
	}
	if !sawAmount {
		t.Fatalf("field deletion not flagged: %+v", rep.Errors)
	}
}

func TestRefuseHugeSpec(t *testing.T) {
	huge := make([]byte, 256*1024+1)
	_, err := CompareYAML(huge, []byte("openapi: 3.0.3\n"))
	if err == nil {
		t.Fatal("expected size error")
	}
}
