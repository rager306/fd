package handlers

import (
	"context"
	"errors"
)

const batchTEISubBatchSize = 32

var errBatchEmbeddingCountMismatch = errors.New("embedding backend returned wrong vector count")

func loadBatchEmbeddings(ctx context.Context, cache EmbeddingCache, embedder Embedder, texts []string, dims, chunkSize int) ([][]float32, error) {
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

func loadBatchEmbeddingChunk(ctx context.Context, cache EmbeddingCache, embedder Embedder, texts []string, vectors [][]float32, chunkStart, chunkEnd, dims int) error {
	missIdx := make([]int, 0, chunkEnd-chunkStart)
	missTexts := make([]string, 0, chunkEnd-chunkStart)

	for i := chunkStart; i < chunkEnd; i++ {
		text := texts[i]
		if emb, ok := cache.GetIfPresent(ctx, text, dims); ok {
			vectors[i] = truncateEmbedding(emb, dims)
			continue
		}
		missIdx = append(missIdx, i)
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
