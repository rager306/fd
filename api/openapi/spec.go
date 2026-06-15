// Package openapi builds fd's OpenAPI 3.1 schema from the implemented HTTP surface.
package openapi

const (
	keyGet         = "get"
	keyPost        = "post"
	keyDescription = "description"
	keyContent     = "content"
	keySchema      = "schema"
	keyType        = "type"
	keyString      = "string"
	keyInteger     = "integer"
	keyNumber      = "number"
	keyObject      = "object"
	keyArray       = "array"
	keyProperties  = "properties"
	keyResponses   = "responses"
	keySummary     = "summary"
)

// Spec returns the fd OpenAPI 3.1 document.
func Spec() map[string]any {
	return m(
		"openapi", "3.1.0",
		"info", m("title", "fd embedding service API", "version", "2.0.0"),
		"paths", paths(),
		"components", components(),
		"security", []map[string]any{{"bearerAuth": []string{}}},
	)
}

func paths() map[string]any {
	return m(
		"/health", getOp("Health", "Deep lifecycle and runtime health", schemaRef("HealthResponse")),
		"/live", getOp("Live", "Cheap process liveness probe", objectSchema()),
		"/ready", getOp("Ready", "Readiness probe after warmup", objectSchema()),
		"/warmup", m(keyGet, operationNoRequest("Warmup status", "Inspect warmup state", objectSchema()), keyPost, operationNoRequest("Trigger warmup", "Trigger model warmup", objectSchema())),
		"/version", getOp("Version", "Build and model version metadata", objectSchema()),
		"/info", getOp("Info", "Service and runtime metadata", objectSchema()),
		"/metrics", textGetOp("Metrics", "Prometheus metrics"),
		"/v1/embeddings", m(keyPost, operation("Create embeddings", "Create OpenAI-compatible embeddings", schemaRef("EmbeddingsRequest"), schemaRef("EmbeddingsResponse"))),
		"/v1/batch", m(keyPost, operation("Create embedding batches", "Create multiple embedding batches", schemaRef("V1BatchRequest"), schemaRef("V1BatchResponse"))),
		"/v1/healthcheck", getOp("V1 healthcheck", "OpenAI-compatible healthcheck alias", schemaRef("HealthResponse")),
		"/v1/traces", getOp("Request traces", "Recent request trace ring buffer", arraySchema(schemaRef("TraceEntry"))),
		"/openapi.json", getOp("OpenAPI schema", "OpenAPI 3.1 JSON schema", objectSchema()),
		"/docs", textGetOp("Swagger UI", "Interactive API documentation"),
	)
}

func getOp(summary, description string, responseSchema map[string]any) map[string]any {
	return m(keyGet, operationNoRequest(summary, description, responseSchema))
}

func textGetOp(summary, description string) map[string]any {
	return m(keyGet, m(
		keySummary, summary,
		keyDescription, description,
		keyResponses, m("200", m(keyDescription, "OK", keyContent, m("text/plain", m(keySchema, stringSchema())))),
	))
}

func operation(summary, description string, requestSchema, responseSchema map[string]any) map[string]any {
	op := operationNoRequest(summary, description, responseSchema)
	op["requestBody"] = m("required", true, keyContent, m("application/json", m(keySchema, requestSchema)))
	return op
}

func operationNoRequest(summary, description string, responseSchema map[string]any) map[string]any {
	return m(
		keySummary, summary,
		keyDescription, description,
		keyResponses, m(
			"200", response("OK", responseSchema),
			"400", errorResponse(),
			"401", errorResponse(),
			"413", errorResponse(),
			"429", errorResponse(),
			"500", errorResponse(),
			"503", errorResponse(),
		),
	)
}

func response(description string, schema map[string]any) map[string]any {
	return m(keyDescription, description, keyContent, m("application/json", m(keySchema, schema)), "headers", responseHeaders())
}

func errorResponse() map[string]any { return response("Error", schemaRef("ErrorResponse")) }

func responseHeaders() map[string]any {
	stringHeader := m(keySchema, stringSchema())
	return m(
		"Server", stringHeader,
		"X-Request-Id", stringHeader,
		"X-Model-Id", stringHeader,
		"X-Dimensions", stringHeader,
		"X-Cache", stringHeader,
		"Retry-After", stringHeader,
		"Cache-Control", stringHeader,
		"ETag", stringHeader,
	)
}

func components() map[string]any {
	return m(
		"securitySchemes", m("bearerAuth", m(keyType, "http", "scheme", "bearer")),
		"schemas", m(
			"EmbeddingsRequest", m(keyType, keyObject, "required", []string{"input"}, keyProperties, m(
				"model", stringSchema(),
				"input", m("oneOf", []map[string]any{stringSchema(), arraySchema(stringSchema())}),
				"dimensions", m(keyType, keyInteger, "enum", []int{512, 1024}),
				"encoding_format", m(keyType, keyString, "enum", []string{"float", "base64"}),
				"user", stringSchema(),
				"priority", m(keyType, keyString, "enum", []string{"low", "normal", "high"}),
			)),
			"EmbeddingsResponse", m(keyType, keyObject, keyProperties, m("object", stringSchema(), "model", stringSchema(), "data", arraySchema(schemaRef("EmbeddingObject")), "usage", schemaRef("Usage"))),
			"EmbeddingObject", m(keyType, keyObject, keyProperties, m("object", stringSchema(), "index", integerSchema(), "dimensions", integerSchema(), "embedding", m("oneOf", []map[string]any{arraySchema(numberSchema()), stringSchema()}))),
			"Usage", m(keyType, keyObject, keyProperties, m("prompt_tokens", integerSchema(), "total_tokens", integerSchema())),
			"V1BatchRequest", m(keyType, keyObject, "required", []string{"batches"}, keyProperties, m("batches", arraySchema(arraySchema(stringSchema())))),
			"V1BatchResponse", m(keyType, keyObject, keyProperties, m("batches", arraySchema(arraySchema(arraySchema(numberSchema()))))),
			"HealthResponse", objectSchema(),
			"TraceEntry", m(keyType, keyObject, keyProperties, m("timestamp", stringSchema(), "latency_ms", integerSchema(), "status", integerSchema(), "model_id", stringSchema(), "request_id", stringSchema(), "path", stringSchema(), "dimensions", integerSchema())),
			"ErrorResponse", m(keyType, keyObject, keyProperties, m("error", m(keyType, keyObject, keyProperties, m("code", stringSchema(), "type", stringSchema(), "param", stringSchema(), "message", stringSchema())))),
		),
	)
}

func schemaRef(name string) map[string]any { return m("$ref", "#/components/schemas/"+name) }
func stringSchema() map[string]any         { return m(keyType, keyString) }
func integerSchema() map[string]any        { return m(keyType, keyInteger) }
func numberSchema() map[string]any         { return m(keyType, keyNumber) }
func objectSchema() map[string]any         { return m(keyType, keyObject, "additionalProperties", true) }
func arraySchema(items map[string]any) map[string]any {
	return m(keyType, keyArray, "items", items)
}

func m(pairs ...any) map[string]any {
	out := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			panic("openapi m key must be string")
		}
		out[key] = pairs[i+1]
	}
	return out
}
