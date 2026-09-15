//go:build analytics_openapi

package buildkite

import (
	"net/http"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"go.yaml.in/yaml/v3"
)

const analyticsOpenAPIURL = "https://api.buildkite.com/v2/analytics/openapi.yaml"

type openAPISchema struct {
	Ref        string                   `yaml:"$ref"`
	AllOf      []openAPISchema          `yaml:"allOf"`
	Properties map[string]openAPISchema `yaml:"properties"`
}

func TestAnalyticsResponseTypesMatchOpenAPI(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(analyticsOpenAPIURL)
	if err != nil {
		t.Fatalf("fetching Analytics OpenAPI document: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("fetching Analytics OpenAPI document: %s", resp.Status)
	}

	var document struct {
		Components struct {
			Schemas map[string]openAPISchema `yaml:"schemas"`
		} `yaml:"components"`
	}
	if err := yaml.NewDecoder(resp.Body).Decode(&document); err != nil {
		t.Fatalf("decoding Analytics OpenAPI document: %v", err)
	}

	responseTypes := map[string]any{
		"BuildExecution":      BuildExecution{},
		"Execution":           FailedExecution{},
		"ExecutionTrace":      ExecutionTrace{},
		"Run":                 TestRun{},
		"Suite":               TestSuite{},
		"Test":                Test{},
		"TestWithMetrics":     TestWithMetrics{},
		"TraceCategoryRollup": TraceCategoryRollup{},
		"TraceSpan":           TraceSpan{},
	}

	for schemaName, responseType := range responseTypes {
		schema, ok := document.Components.Schemas[schemaName]
		if !ok {
			t.Errorf("OpenAPI response schema %q no longer exists", schemaName)
			continue
		}

		want := schemaProperties(schema, document.Components.Schemas)
		got := jsonFields(reflect.TypeOf(responseType))
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%T JSON fields differ from OpenAPI schema %s\n  Go:      %s\n  OpenAPI: %s", responseType, schemaName, strings.Join(got, ", "), strings.Join(want, ", "))
		}
	}
}

func schemaProperties(schema openAPISchema, schemas map[string]openAPISchema) []string {
	properties := make(map[string]bool)
	var addProperties func(openAPISchema)
	addProperties = func(schema openAPISchema) {
		if schema.Ref != "" {
			name := strings.TrimPrefix(schema.Ref, "#/components/schemas/")
			if referenced, ok := schemas[name]; ok {
				addProperties(referenced)
			}
		}
		for _, part := range schema.AllOf {
			addProperties(part)
		}
		for name := range schema.Properties {
			properties[name] = true
		}
	}
	addProperties(schema)

	return sortedKeys(properties)
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
