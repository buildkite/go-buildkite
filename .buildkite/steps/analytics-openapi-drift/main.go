package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	buildkite "github.com/buildkite/go-buildkite/v5"
)

const openAPIURL = "https://api.buildkite.com/v2/analytics/openapi.yaml"

type schema struct {
	properties map[string]bool
	refs       []string
}

func main() {
	schemas, err := fetchSchemas()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	responseTypes := []struct {
		schema string
		value  any
	}{
		{"BuildExecution", buildkite.BuildExecution{}},
		{"Execution", buildkite.FailedExecution{}},
		{"ExecutionTrace", buildkite.ExecutionTrace{}},
		{"Run", buildkite.TestRun{}},
		{"Suite", buildkite.TestSuite{}},
		{"Test", buildkite.Test{}},
		{"TestWithMetrics", buildkite.TestWithMetrics{}},
		{"TraceCategoryRollup", buildkite.TraceCategoryRollup{}},
		{"TraceSpan", buildkite.TraceSpan{}},
	}

	drift := false
	for _, responseType := range responseTypes {
		want, err := schemaProperties(responseType.schema, schemas, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			drift = true
			continue
		}
		got := jsonFields(reflect.TypeOf(responseType.value))
		if !reflect.DeepEqual(got, want) {
			fmt.Fprintf(os.Stderr, "%T JSON fields differ from OpenAPI schema %s\n  Go:      %s\n  OpenAPI: %s\n", responseType.value, responseType.schema, strings.Join(got, ", "), strings.Join(want, ", "))
			drift = true
		}
	}

	if drift {
		os.Exit(1)
	}
}

func fetchSchemas() (map[string]schema, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(openAPIURL)
	if err != nil {
		return nil, fmt.Errorf("fetching Analytics OpenAPI document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching Analytics OpenAPI document: %s", resp.Status)
	}

	schemas := make(map[string]schema)
	var current string
	inSchemas := false
	propertiesIndent := -1
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if line == "  schemas:" {
			inSchemas = true
			continue
		}
		if !inSchemas {
			continue
		}
		if trimmed != "" && indent <= 2 {
			break
		}
		if indent == 4 && strings.HasSuffix(trimmed, ":") {
			current = strings.TrimSuffix(trimmed, ":")
			schemas[current] = schema{properties: make(map[string]bool)}
			propertiesIndent = -1
			continue
		}
		if current == "" {
			continue
		}
		if trimmed == "properties:" {
			propertiesIndent = indent
			continue
		}
		if propertiesIndent >= 0 && indent <= propertiesIndent && trimmed != "" {
			propertiesIndent = -1
		}
		if propertiesIndent >= 0 && indent == propertiesIndent+2 && strings.HasSuffix(trimmed, ":") {
			name := strings.Trim(strings.TrimSuffix(trimmed, ":"), `"'`)
			s := schemas[current]
			s.properties[name] = true
			schemas[current] = s
		}
		if indent == 6 && strings.HasPrefix(trimmed, `- "$ref":`) {
			ref := strings.TrimSpace(strings.TrimPrefix(trimmed, `- "$ref":`))
			ref = strings.Trim(ref, `"'`)
			s := schemas[current]
			s.refs = append(s.refs, strings.TrimPrefix(ref, "#/components/schemas/"))
			schemas[current] = s
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading Analytics OpenAPI document: %w", err)
	}
	return schemas, nil
}

func schemaProperties(name string, schemas map[string]schema, visiting map[string]bool) ([]string, error) {
	s, ok := schemas[name]
	if !ok {
		return nil, fmt.Errorf("OpenAPI response schema %q no longer exists", name)
	}
	if visiting == nil {
		visiting = make(map[string]bool)
	}
	if visiting[name] {
		return nil, fmt.Errorf("OpenAPI response schema %q contains a reference cycle", name)
	}
	visiting[name] = true
	defer delete(visiting, name)

	properties := make(map[string]bool)
	for property := range s.properties {
		properties[property] = true
	}
	for _, ref := range s.refs {
		referenced, err := schemaProperties(ref, schemas, visiting)
		if err != nil {
			return nil, err
		}
		for _, property := range referenced {
			properties[property] = true
		}
	}
	return sortedKeys(properties), nil
}

func jsonFields(typ reflect.Type) []string {
	fields := make(map[string]bool)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if field.Anonymous && name == "" {
			for _, embedded := range jsonFields(field.Type) {
				fields[embedded] = true
			}
			continue
		}
		if name != "" && name != "-" {
			fields[name] = true
		}
	}
	return sortedKeys(fields)
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
