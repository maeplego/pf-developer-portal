package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const maxSpecBytes = 256 * 1024

// API is one hand-placed OpenAPI document.
type API struct {
	Slug     string
	Title    string
	Version  string
	Summary  string
	Source   string
	FileName string
	Raw      []byte
	Doc      map[string]any
	Paths    []PathOp
}

type PathOp struct {
	Method      string
	Path        string
	OperationID string
	Summary     string
	Tag         string
	Op          map[string]any
}

type Catalog struct {
	APIs []API
	by   map[string]*API
}

func LoadDir(dir string) (*Catalog, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	c := &Catalog{by: map[string]*API{}}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		low := strings.ToLower(name)
		if !strings.HasSuffix(low, ".yaml") && !strings.HasSuffix(low, ".yml") {
			continue
		}
		api, err := LoadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if _, exists := c.by[api.Slug]; exists {
			return nil, fmt.Errorf("duplicate slug %q", api.Slug)
		}
		c.APIs = append(c.APIs, api)
		c.by[api.Slug] = &c.APIs[len(c.APIs)-1]
	}
	if len(c.APIs) == 0 {
		return nil, fmt.Errorf("no OpenAPI yaml in %s", dir)
	}
	sort.Slice(c.APIs, func(i, j int) bool { return c.APIs[i].Slug < c.APIs[j].Slug })
	c.by = map[string]*API{}
	for i := range c.APIs {
		c.by[c.APIs[i].Slug] = &c.APIs[i]
	}
	return c, nil
}

func LoadFile(path string) (API, error) {
	st, err := os.Stat(path)
	if err != nil {
		return API{}, err
	}
	if st.Size() > maxSpecBytes {
		return API{}, fmt.Errorf("spec larger than %d bytes", maxSpecBytes)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return API{}, err
	}
	if len(raw) > maxSpecBytes {
		return API{}, fmt.Errorf("spec larger than %d bytes", maxSpecBytes)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return API{}, err
	}
	info := asMap(doc["info"])
	title := asString(info["title"])
	if title == "" {
		return API{}, fmt.Errorf("info.title required")
	}
	portal := asMap(info["x-portal"])
	slug := asString(portal["slug"])
	if slug == "" {
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		slug = base
	}
	api := API{
		Slug:     slug,
		Title:    title,
		Version:  asString(info["version"]),
		Summary:  asString(portal["summary"]),
		Source:   asString(portal["source"]),
		FileName: filepath.Base(path),
		Raw:      raw,
		Doc:      doc,
		Paths:    collectOps(doc),
	}
	return api, nil
}

func (c *Catalog) Get(slug string) (*API, bool) {
	a, ok := c.by[slug]
	return a, ok
}

func collectOps(doc map[string]any) []PathOp {
	paths := asMap(doc["paths"])
	keys := make([]string, 0, len(paths))
	for k := range paths {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	methods := []string{"get", "post", "put", "patch", "delete", "head"}
	var out []PathOp
	for _, p := range keys {
		item := asMap(paths[p])
		for _, m := range methods {
			op := asMap(item[m])
			if len(op) == 0 {
				continue
			}
			tag := ""
			if tags, ok := op["tags"].([]any); ok && len(tags) > 0 {
				tag = asString(tags[0])
			}
			out = append(out, PathOp{
				Method:      strings.ToUpper(m),
				Path:        p,
				OperationID: asString(op["operationId"]),
				Summary:     asString(op["summary"]),
				Tag:         tag,
				Op:          op,
			})
		}
	}
	return out
}

func asMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}
