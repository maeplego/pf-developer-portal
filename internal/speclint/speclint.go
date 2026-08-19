package speclint

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Patterns that must not appear in portal examples (mask-worthy, not exploits).
var forbidden = []*regexp.Regexp{
	regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
	regexp.MustCompile(`ghp_[A-Za-z0-9]{20,}`),
	regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`),
}

func ScanDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var problems []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		low := strings.ToLower(name)
		if !strings.HasSuffix(low, ".yaml") && !strings.HasSuffix(low, ".yml") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		text := string(b)
		for _, re := range forbidden {
			if re.MatchString(text) {
				problems = append(problems, fmt.Sprintf("%s matches %s", name, re.String()))
			}
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("example secret lint: %s", strings.Join(problems, "; "))
	}
	return nil
}
