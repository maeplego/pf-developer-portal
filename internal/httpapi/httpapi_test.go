package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/portfolio/pf-developer-portal/internal/catalog"
)

func testServer(t *testing.T) http.Handler {
	t.Helper()
	c, err := catalog.LoadDir(filepath.Join("..", "..", "specs"))
	if err != nil {
		t.Fatal(err)
	}
	return New(c)
}

func TestHealthAndCatalog(t *testing.T) {
	h := testServer(t)
	for _, p := range []string{"/health", "/ready"} {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, p, nil))
		if rr.Code != 200 {
			t.Fatalf("%s %d", p, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), `"ok":true`) {
			t.Fatalf("%s body %s", p, rr.Body.String())
		}
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/catalog", nil))
	if rr.Code != 200 {
		t.Fatal(rr.Code)
	}
	var body struct {
		APIs []struct {
			Slug string `json:"slug"`
		}
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.APIs) != 3 {
		t.Fatalf("%+v", body)
	}
}

func TestDocsRender(t *testing.T) {
	h := testServer(t)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != 200 || !strings.Contains(rr.Body.String(), "Payments API") {
		t.Fatalf("%d %s", rr.Code, rr.Body.String())
	}
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/docs/payments", nil))
	if rr.Code != 200 {
		t.Fatal(rr.Code)
	}
	html := rr.Body.String()
	for _, want := range []string{"POST", "/v1/charges", "Idempotency-Key", "Try the mock"} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q", want)
		}
	}
}

func TestAPIDiffBreaking(t *testing.T) {
	h := testServer(t)
	base, err := os.ReadFile(filepath.Join("..", "..", "testdata", "openapi", "base.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	brk, err := os.ReadFile(filepath.Join("..", "..", "testdata", "openapi", "breaking.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(map[string]string{"base": string(base), "revision": string(brk)})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/diff", strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rr, req)
	if rr.Code != 409 {
		t.Fatalf("want 409 got %d %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "response-property-removed") {
		t.Fatalf("%s", rr.Body.String())
	}
}

func TestMockExampleAndValidation(t *testing.T) {
	h := testServer(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/mock/payments/v1/charges", strings.NewReader(`{"amountMinor":1299,"currency":"JPY","source":"tok_demo_visa"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "demo-charge-001")
	h.ServeHTTP(rr, req)
	if rr.Code != 201 {
		t.Fatalf("create %d %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ch_demo_1") {
		t.Fatalf("example %s", rr.Body.String())
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/mock/payments/v1/charges", strings.NewReader(`{"currency":"JPY"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "demo-charge-002")
	h.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Fatalf("want 400 got %d %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/mock/payments/v1/charges", strings.NewReader(`{"amountMinor":1,"currency":"JPY","source":"tok"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Fatalf("missing header %d %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/mock/payments/v1/charges/ch_demo_1", nil))
	if rr.Code != 200 {
		t.Fatalf("get %d %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/mock/payments/nope", nil))
	if rr.Code != 404 {
		t.Fatalf("404 %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/mock/commerce-catalog/v1/products", nil))
	if rr.Code != 200 {
		t.Fatalf("catalog %d %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/specs/content-blog", nil))
	if rr.Code != 200 {
		t.Fatal(rr.Code)
	}
	raw, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(raw), "hello-portal") {
		t.Fatalf("%s", raw)
	}
}
