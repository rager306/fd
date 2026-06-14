package handlers

const (
	// maxBatchInputChars mirrors the /v1/embeddings per-input validation
	// bound. Batch endpoints use distinct JSON shapes, so they validate this
	// in handlers before cache/TEI work rather than reusing the embeddings
	// request middleware directly.
	maxBatchInputChars = 2048

	// maxLegacyBatchInputs keeps the legacy /embeddings/batch endpoint within
	// the same request input-count envelope as /v1/embeddings.
	maxLegacyBatchInputs = 128
)
