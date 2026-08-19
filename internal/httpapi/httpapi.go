package httpapi

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/portfolio/pf-developer-portal/internal/catalog"
	"github.com/portfolio/pf-developer-portal/internal/mock"
)

type Server struct {
	cat  *catalog.Catalog
	page *template.Template
}

func New(cat *catalog.Catalog) http.Handler {
	s := &Server{cat: cat, page: mustParse()}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /ready", s.health)
	mux.HandleFunc("GET /api/catalog", s.apiCatalog)
	mux.HandleFunc("GET /api/specs/{slug}", s.apiSpec)
	mux.HandleFunc("GET /docs/{slug}", s.docs)
	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("/mock/{slug}", s.mock)
	mux.HandleFunc("/mock/{slug}/{path...}", s.mock)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func (s *Server) apiCatalog(w http.ResponseWriter, _ *http.Request) {
	type item struct {
		Slug    string `json:"slug"`
		Title   string `json:"title"`
		Version string `json:"version"`
		Summary string `json:"summary"`
		Source  string `json:"source"`
		Ops     int    `json:"operations"`
	}
	out := make([]item, 0, len(s.cat.APIs))
	for _, a := range s.cat.APIs {
		out = append(out, item{a.Slug, a.Title, a.Version, a.Summary, a.Source, len(a.Paths)})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"apis": out})
}

func (s *Server) apiSpec(w http.ResponseWriter, r *http.Request) {
	api, ok := s.cat.Get(r.PathValue("slug"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write(api.Raw)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	s.render(w, pageData{
		Title: "API catalog",
		Home:  true,
		APIs:  s.cat.APIs,
	})
}

func (s *Server) docs(w http.ResponseWriter, r *http.Request) {
	api, ok := s.cat.Get(r.PathValue("slug"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	ops := make([]opView, 0, len(api.Paths))
	for _, p := range api.Paths {
		ops = append(ops, toOpView(api, p))
	}
	s.render(w, pageData{
		Title: api.Title,
		API:   api,
		APIs:  s.cat.APIs,
		Ops:   ops,
		Mock:  "/mock/" + api.Slug,
	})
}

func (s *Server) mock(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	path := r.PathValue("path")
	api, ok := s.cat.Get(slug)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":"not_found","message":"unknown api"}}`))
		return
	}
	r2 := r.Clone(r.Context())
	if path == "" {
		r2.URL.Path = "/"
	} else {
		r2.URL.Path = "/" + path
	}
	mock.Server{API: api}.ServeHTTP(w, r2)
}

func (s *Server) render(w http.ResponseWriter, data pageData) {
	var buf bytes.Buffer
	if err := s.page.Execute(&buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

type pageData struct {
	Title string
	Home  bool
	APIs  []catalog.API
	API   *catalog.API
	Ops   []opView
	Mock  string
}

type opView struct {
	ID          string
	Method      string
	Path        string
	Summary     string
	Description string
	Tag         string
	Params      []paramView
	BodyExample string
	RespExample string
	RespStatus  string
}

type paramView struct {
	In       string
	Name     string
	Required bool
	Type     string
}

func toOpView(api *catalog.API, p catalog.PathOp) opView {
	v := opView{
		ID:          p.Method + "-" + strings.ReplaceAll(p.Path, "/", "-"),
		Method:      p.Method,
		Path:        p.Path,
		Summary:     p.Summary,
		Description: asString(p.Op["description"]),
		Tag:         p.Tag,
	}
	raw, _ := p.Op["parameters"].([]any)
	for _, item := range raw {
		m, _ := item.(map[string]any)
		if m == nil {
			continue
		}
		schema := asMap(m["schema"])
		req, _ := m["required"].(bool)
		v.Params = append(v.Params, paramView{
			In:       asString(m["in"]),
			Name:     asString(m["name"]),
			Required: req,
			Type:     asString(schema["type"]),
		})
	}
	rb := asMap(p.Op["requestBody"])
	media := asMap(asMap(rb["content"])["application/json"])
	if ex, ok := media["example"]; ok {
		v.BodyExample = pretty(ex)
	}
	resps := asMap(p.Op["responses"])
	for _, code := range []string{"200", "201", "202"} {
		if raw, ok := resps[code]; ok {
			v.RespStatus = code
			v.RespExample = pretty(exampleJSON(api.Doc, raw))
			break
		}
	}
	return v
}

func exampleJSON(doc map[string]any, raw any) any {
	m, _ := raw.(map[string]any)
	if m == nil {
		return nil
	}
	if ref, _ := m["$ref"].(string); ref != "" {
		cur := any(doc)
		for _, p := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
			mm, _ := cur.(map[string]any)
			cur = mm[p]
		}
		m, _ = cur.(map[string]any)
	}
	media := asMap(asMap(m["content"])["application/json"])
	if ex, ok := media["example"]; ok {
		return ex
	}
	return nil
}

func pretty(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func asMap(v any) map[string]any {
	m, _ := v.(map[string]any)
	if m == nil {
		return map[string]any{}
	}
	return m
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}
