package embed

// ONNXEmbedderOptions configures the ONNX runtime: manifest path,
// runtime shared library path, tokenizer path, and maximum sequence length.
type ONNXEmbedderOptions struct {
	ManifestPath      string
	SharedLibraryPath string
	TokenizerPath     string
	MaxSequenceLength int
}
