package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/portfolio/pf-developer-portal/internal/catalog"
)

const maxBody = 64 * 1024

type Server struct {
	API *catalog.API
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Idempotency-Key")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	op, params, ok := findOp(s.API, r.Method, r.URL.Path)
	if !ok {
		writeErr(w, http.StatusNotFound, "not_found", "path or method is not in this spec")
		return
	}
	if err := validateParams(op, r, params); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	body, err := readJSON(r, op)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if err := validateBody(s.API.Doc, op, body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	status, payload := pickExample(s.API.Doc, op)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(payload)
}

func findOp(api *catalog.API, method, reqPath string) (map[string]any, map[string]string, bool) {
	method = strings.ToUpper(method)
	for _, p := range api.Paths {
		if p.Method != method {
			continue
		}
		params, ok := matchPath(p.Path, reqPath)
		if ok {
			return p.Op, params, true
		}
	}
	return nil, nil, false
}

func matchPath(specPath, reqPath string) (map[string]string, bool) {
	sp := splitPath(specPath)
	rp := splitPath(reqPath)
	if len(sp) != len(rp) {
		return nil, false
	}
	params := map[string]string{}
	for i := range sp {
		if strings.HasPrefix(sp[i], "{") && strings.HasSuffix(sp[i], "}") {
			params[sp[i][1:len(sp[i])-1]] = rp[i]
			continue
		}
		if sp[i] != rp[i] {
			return nil, false
		}
	}
	return params, true
}

func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

func validateParams(op map[string]any, r *http.Request, pathParams map[string]string) error {
	raw, _ := op["parameters"].([]any)
	for _, item := range raw {
		p, _ := item.(map[string]any)
		if p == nil {
			continue
		}
		required, _ := p["required"].(bool)
		if !required {
			continue
		}
		name := asString(p["name"])
		switch asString(p["in"]) {
		case "path":
			if pathParams[name] == "" {
				return fmt.Errorf("path parameter %s is required", name)
			}
		case "query":
			if r.URL.Query().Get(name) == "" {
				return fmt.Errorf("query parameter %s is required", name)
			}
		case "header":
			if r.Header.Get(name) == "" {
				return fmt.Errorf("header %s is required", name)
			}
		}
	}
	return nil
}

func readJSON(r *http.Request, op map[string]any) (any, error) {
	rb := asMap(op["requestBody"])
	if len(rb) == 0 {
		return nil, nil
	}
	required, _ := rb["required"].(bool)
	r.Body = http.MaxBytesReader(nil, r.Body, maxBody)
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("body too large or unreadable")
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		if required {
			return nil, fmt.Errorf("request body is required")
		}
		return nil, nil
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("body must be JSON")
	}
	return v, nil
}

func validateBody(doc map[string]any, op map[string]any, body any) error {
	rb := asMap(op["requestBody"])
	if len(rb) == 0 {
		return nil
	}
	required, _ := rb["required"].(bool)
	if body == nil {
		if required {
			return fmt.Errorf("request body is required")
		}
		return nil
	}
	content := asMap(rb["content"])
	media := asMap(content["application/json"])
	schema := resolveSchema(doc, asMap(media["schema"]))
	return checkSchema(doc, schema, body, "")
}

func checkSchema(doc, schema map[string]any, val any, path string) error {
	schema = resolveSchema(doc, schema)
	if len(schema) == 0 {
		return nil
	}
	typ := asString(schema["type"])
	at := func(name string) string {
		if path == "" {
			return name
		}
		if name == "" {
			return path
		}
		return path + "." + name
	}
	switch typ {
	case "object", "":
		m, ok := val.(map[string]any)
		if !ok && typ == "object" {
			return fmt.Errorf("%s must be an object", orRoot(path))
		}
		if !ok {
			return nil
		}
		req, _ := schema["required"].([]any)
		for _, r := range req {
			key := asString(r)
			if _, exists := m[key]; !exists {
				return fmt.Errorf("%s is required", at(key))
			}
		}
		props := asMap(schema["properties"])
		for k, pv := range m {
			if _, exists := props[k]; !exists && len(props) > 0 {
				continue
			}
			if err := checkSchema(doc, asMap(props[k]), pv, at(k)); err != nil {
				return err
			}
		}
	case "string":
		if _, ok := val.(string); !ok {
			return fmt.Errorf("%s must be a string", orRoot(path))
		}
	case "integer":
		switch val.(type) {
		case float64:
			n, _ := val.(float64)
			if n != float64(int64(n)) {
				return fmt.Errorf("%s must be an integer", orRoot(path))
			}
		case json.Number, int, int64:
		default:
			return fmt.Errorf("%s must be an integer", orRoot(path))
		}
	case "number":
		switch val.(type) {
		case float64, json.Number, int, int64:
		default:
			return fmt.Errorf("%s must be a number", orRoot(path))
		}
	case "boolean":
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("%s must be a boolean", orRoot(path))
		}
	case "array":
		arr, ok := val.([]any)
		if !ok {
			return fmt.Errorf("%s must be an array", orRoot(path))
		}
		for i, item := range arr {
			if err := checkSchema(doc, asMap(schema["items"]), item, fmt.Sprintf("%s[%d]", orRoot(path), i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func pickExample(doc map[string]any, op map[string]any) (int, any) {
	resps := asMap(op["responses"])
	for _, code := range []string{"200", "201", "202"} {
		if raw, ok := resps[code]; ok {
			status := 200
			fmt.Sscanf(code, "%d", &status)
			body := exampleFromResponse(doc, raw)
			if body == nil {
				body = map[string]any{"ok": true}
			}
			return status, body
		}
	}
	return http.StatusOK, map[string]any{"ok": true}
}

func exampleFromResponse(doc map[string]any, raw any) any {
	m := asMap(raw)
	if ref := asString(m["$ref"]); ref != "" {
		m = resolveRef(doc, ref)
	}
	content := asMap(m["content"])
	media := asMap(content["application/json"])
	if ex, ok := media["example"]; ok {
		return ex
	}
	examples := asMap(media["examples"])
	for _, v := range examples {
		em := asMap(v)
		if ex, ok := em["value"]; ok {
			return ex
		}
	}
	return nil
}

func resolveSchema(doc map[string]any, schema map[string]any) map[string]any {
	if ref := asString(schema["$ref"]); ref != "" && doc != nil {
		return resolveSchema(doc, resolveRef(doc, ref))
	}
	return schema
}

func resolveRef(doc map[string]any, ref string) map[string]any {
	if !strings.HasPrefix(ref, "#/") {
		return map[string]any{}
	}
	cur := any(doc)
	for _, p := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		m, ok := cur.(map[string]any)
		if !ok {
			return map[string]any{}
		}
		cur = m[p]
	}
	m, _ := cur.(map[string]any)
	if m == nil {
		return map[string]any{}
	}
	return m
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": code, "message": msg},
	})
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

func orRoot(path string) string {
	if path == "" {
		return "value"
	}
	return path
}
