package specbreak

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ERR-level subset used as the local oasdiff-style gate (path removed,
// required request field added, response property removed, type change).
// GitHub Actions still runs upstream oasdiff-action on the same fixtures.

type Finding struct {
	ID      string `json:"id"`
	Method  string `json:"method,omitempty"`
	Path    string `json:"path,omitempty"`
	Message string `json:"message"`
}

type Report struct {
	Errors []Finding `json:"errors"`
}

func CompareYAML(baseRaw, revRaw []byte) (Report, error) {
	base, err := parse(baseRaw)
	if err != nil {
		return Report{}, fmt.Errorf("base: %w", err)
	}
	rev, err := parse(revRaw)
	if err != nil {
		return Report{}, fmt.Errorf("revision: %w", err)
	}
	return compare(base, rev), nil
}

func (r Report) HasERR() bool {
	return len(r.Errors) > 0
}

type doc struct {
	raw   map[string]any
	ops   map[string]map[string]any // METHOD path -> operation
	comps map[string]any
}

func parse(raw []byte) (*doc, error) {
	if len(raw) > 256*1024 {
		return nil, fmt.Errorf("spec larger than 256KiB")
	}
	var m map[string]any
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	if m == nil {
		return nil, fmt.Errorf("empty yaml")
	}
	d := &doc{raw: m, ops: map[string]map[string]any{}, comps: asMap(asMap(m["components"])["schemas"])}
	paths := asMap(m["paths"])
	for p, rawOp := range paths {
		methods := asMap(rawOp)
		for method, op := range methods {
			ml := strings.ToUpper(method)
			switch ml {
			case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
			default:
				continue
			}
			om := asMap(op)
			d.ops[ml+" "+p] = om
		}
	}
	return d, nil
}

func compare(base, rev *doc) Report {
	var errs []Finding
	var keys []string
	for k := range base.ops {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		parts := strings.SplitN(k, " ", 2)
		method, path := parts[0], parts[1]
		rop, ok := rev.ops[k]
		if !ok {
			errs = append(errs, Finding{
				ID:      "api-path-removed",
				Method:  method,
				Path:    path,
				Message: "operation removed",
			})
			continue
		}
		bop := base.ops[k]
		errs = append(errs, requestRequiredAdded(base, rev, method, path, bop, rop)...)
		errs = append(errs, responsePropsRemoved(base, rev, method, path, bop, rop)...)
		errs = append(errs, paramRequiredAdded(method, path, bop, rop)...)
	}
	return Report{Errors: errs}
}

func requestRequiredAdded(base, rev *doc, method, path string, bop, rop map[string]any) []Finding {
	bs := schemaFromRequest(base, bop)
	rs := schemaFromRequest(rev, rop)
	if bs == nil || rs == nil {
		return nil
	}
	bold := requiredSet(bs)
	rnew := requiredSet(rs)
	var out []Finding
	for name := range rnew {
		if !bold[name] {
			out = append(out, Finding{
				ID:      "request-property-became-required",
				Method:  method,
				Path:    path,
				Message: "required request property added: " + name,
			})
		}
	}
	bprops := propsOf(bs)
	rprops := propsOf(rs)
	for name, bsch := range bprops {
		rsch, ok := rprops[name]
		if !ok {
			continue
		}
		bt := asString(bsch["type"])
		rt := asString(rsch["type"])
		if bt != "" && rt != "" && bt != rt {
			out = append(out, Finding{
				ID:      "request-property-type-changed",
				Method:  method,
				Path:    path,
				Message: fmt.Sprintf("request property %s type %s -> %s", name, bt, rt),
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Message < out[j].Message })
	return out
}

func responsePropsRemoved(base, rev *doc, method, path string, bop, rop map[string]any) []Finding {
	var out []Finding
	bres := asMap(bop["responses"])
	rres := asMap(rop["responses"])
	for code, raw := range bres {
		if !isSuccess(code) {
			continue
		}
		bs := schemaFromResponse(base, asMap(raw))
		rs := schemaFromResponse(rev, asMap(rres[code]))
		if bs == nil {
			continue
		}
		if rres[code] == nil {
			out = append(out, Finding{
				ID:      "response-success-status-removed",
				Method:  method,
				Path:    path,
				Message: "success response " + code + " removed",
			})
			continue
		}
		if rs == nil {
			continue
		}
		bprops := propsOf(bs)
		rprops := propsOf(rs)
		var names []string
		for n := range bprops {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, name := range names {
			if _, ok := rprops[name]; !ok {
				out = append(out, Finding{
					ID:      "response-property-removed",
					Method:  method,
					Path:    path,
					Message: "response " + code + " property removed: " + name,
				})
			}
		}
	}
	return out
}

func paramRequiredAdded(method, path string, bop, rop map[string]any) []Finding {
	b := paramsByKey(bop)
	r := paramsByKey(rop)
	var out []Finding
	for k, rp := range r {
		req, _ := rp["required"].(bool)
		if !req {
			continue
		}
		bp, ok := b[k]
		if !ok {
			out = append(out, Finding{
				ID:      "request-parameter-became-required",
				Method:  method,
				Path:    path,
				Message: "required parameter added: " + k,
			})
			continue
		}
		was, _ := bp["required"].(bool)
		if !was {
			out = append(out, Finding{
				ID:      "request-parameter-became-required",
				Method:  method,
				Path:    path,
				Message: "parameter became required: " + k,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Message < out[j].Message })
	return out
}

func paramsByKey(op map[string]any) map[string]map[string]any {
	out := map[string]map[string]any{}
	raw, _ := op["parameters"].([]any)
	for _, item := range raw {
		m := asMap(item)
		in := asString(m["in"])
		name := asString(m["name"])
		if in == "" || name == "" {
			continue
		}
		out[in+":"+name] = m
	}
	return out
}

func schemaFromRequest(d *doc, op map[string]any) map[string]any {
	rb := asMap(op["requestBody"])
	media := asMap(asMap(rb["content"])["application/json"])
	return resolveSchema(d, media["schema"])
}

func schemaFromResponse(d *doc, resp map[string]any) map[string]any {
	if ref := asString(resp["$ref"]); ref != "" {
		resp = resolveRef(d.raw, ref)
	}
	media := asMap(asMap(resp["content"])["application/json"])
	return resolveSchema(d, media["schema"])
}

func resolveSchema(d *doc, raw any) map[string]any {
	m := asMap(raw)
	if ref := asString(m["$ref"]); ref != "" {
		m = resolveRef(d.raw, ref)
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func resolveRef(doc map[string]any, ref string) map[string]any {
	if !strings.HasPrefix(ref, "#/") {
		return nil
	}
	cur := any(doc)
	for _, p := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		mm, _ := cur.(map[string]any)
		if mm == nil {
			return nil
		}
		cur = mm[p]
	}
	return asMap(cur)
}

func requiredSet(schema map[string]any) map[string]bool {
	out := map[string]bool{}
	raw, _ := schema["required"].([]any)
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out[s] = true
		}
	}
	return out
}

func propsOf(schema map[string]any) map[string]map[string]any {
	out := map[string]map[string]any{}
	for k, v := range asMap(schema["properties"]) {
		out[k] = asMap(v)
	}
	return out
}

func isSuccess(code string) bool {
	return strings.HasPrefix(code, "2")
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
