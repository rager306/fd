package openapi

import (
	"encoding/json"
	"testing"
)

func TestSpecIsOpenAPI31JSONWithExpectedPaths(t *testing.T) {
	spec := Spec()
	if got := spec["openapi"]; got != "3.1.0" {
		t.Fatalf("openapi = %v, want 3.1.0", got)
	}
	if _, err := json.Marshal(spec); err != nil {
		t.Fatalf("spec is not JSON-marshalable: %v", err)
	}
	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatalf("paths missing or wrong type")
	}
	for _, path := range []string{"/health", "/live", "/ready", "/warmup", "/version", "/info", "/metrics", "/v1/embeddings", "/v1/batch", "/v1/healthcheck", "/v1/traces", "/openapi.json", "/docs"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("path %s missing", path)
		}
	}
	components, ok := spec["components"].(map[string]any)
	if !ok {
		t.Fatalf("components missing or wrong type")
	}
	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		t.Fatalf("components.schemas missing or wrong type")
	}
	for _, name := range []string{"EmbeddingsRequest", "EmbeddingsResponse", "V1BatchRequest", "V1BatchResponse", "ErrorResponse", "TraceEntry"} {
		if _, ok := schemas[name]; !ok {
			t.Fatalf("schema %s missing", name)
		}
	}
}
