package handlers

import (
	"context"
	"errors"

	"fd-api/embed"
)

const batchTEISubBatchSize = 32

var errBatchEmbeddingCountMismatch = errors.New("embedding backend returned wrong vector count")

func loadBatchEmbeddings(ctx context.Context, cache EmbeddingCache, embedder embed.Embedder, texts []string, dims, chunkSize int) ([][]float32, error) {
	if chunkSize <= 0 {
		chunkSize = batchTEISubBatchSize
	}

	vectors := make([][]float32, len(texts))
	for chunkStart := 0; chunkStart < len(texts); chunkStart += chunkSize {
		chunkEnd := min(chunkStart+chunkSize, len(texts))
		if err := loadBatchEmbeddingChunk(ctx, cache, embedder, texts, vectors, chunkStart, chunkEnd, dims); err != nil {
			return nil, err
		}
	}
	return vectors, nil
}

func loadBatchEmbeddingChunk(ctx context.Context, cache EmbeddingCache, embedder embed.Embedder, texts []string, vectors [][]float32, chunkStart, chunkEnd, dims int) error {
	chunk := texts[chunkStart:chunkEnd]
	missIdx := make([]int, 0, len(chunk))
	missTexts := make([]string, 0, len(chunk))

	// One lookup round trip per chunk: GetManyIfPresent resolves the L2 tier
	// with a single batched read, where GetIfPresent per text issued one Redis
	// call per text (up to 32 per chunk, 100 chunks for a 3200-input batch).
	hits := cache.GetManyIfPresent(ctx, chunk, dims)
	for j, text := range chunk {
		if emb, ok := hits[j]; ok {
			vectors[chunkStart+j] = truncateEmbedding(emb, dims)
			continue
		}
		missIdx = append(missIdx, chunkStart+j)
		missTexts = append(missTexts, text)
	}
	if len(missTexts) == 0 {
		return nil
	}

	embs, err := embedder.Embed(ctx, missTexts)
	if err != nil {
		return err
	}
	if len(embs) != len(missTexts) {
		return errBatchEmbeddingCountMismatch
	}
	for j, originalIndex := range missIdx {
		emb := embs[j]
		cache.Set(ctx, missTexts[j], dims, emb)
		vectors[originalIndex] = truncateEmbedding(emb, dims)
	}
	return nil
}
