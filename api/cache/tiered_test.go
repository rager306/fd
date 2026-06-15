package cache

import (
	"bytes"
	"context"
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func floatsToBytes(t *testing.T, emb []float32) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	require.NoError(t, binary.Write(buf, binary.LittleEndian, emb))
	return buf.Bytes()
}

func TestTieredCacheGetManyIfPresentUsesLocalHitsByIndex(t *testing.T) {
	ctx := context.Background()
	local := NewLocalCache(100, time.Minute)
	defer func() { require.NoError(t, local.Close()) }()
	tc := NewTieredCacheWithLogger(local, nil, time.Minute, newDiscardLogger())

	vecA := make([]float32, 1024)
	vecA[0] = 1
	vecC := make([]float32, 1024)
	vecC[0] = 3
	dataA, err := marshalEmbedding(vecA, 1024)
	require.NoError(t, err)
	dataC, err := marshalEmbedding(vecC, 1024)
	require.NoError(t, err)
	local.Set(ctx, localCacheKey("a", 1024), dataA, time.Minute)
	local.Set(ctx, localCacheKey("c", 1024), dataC, time.Minute)

	hits := tc.GetManyIfPresent(ctx, []string{"a", "b", "c"}, 1024)
	require.Len(t, hits, 2)
	assert.Equal(t, float32(1), hits[0][0])
	_, ok := hits[1]
	assert.False(t, ok, "middle miss should be omitted")
	assert.Equal(t, float32(3), hits[2][0])
}

func TestRedisBulkBytes(t *testing.T) {
	data := []byte{1, 2, 3}
	got, ok := redisBulkBytes(string(data))
	require.True(t, ok)
	assert.Equal(t, data, got)
	got, ok = redisBulkBytes(data)
	require.True(t, ok)
	assert.Equal(t, data, got)
	_, ok = redisBulkBytes(nil)
	assert.False(t, ok)
	_, ok = redisBulkBytes(123)
	assert.False(t, ok)
}

func TestFloatsToBytes_RoundTrip(t *testing.T) {
	original := []float32{1.5, 2.5, -0.5, 1024.0}
	data := floatsToBytes(t, original)

	var restored []float32
	for i := 0; i < len(original); i++ {
		var v float32
		require.NoError(t, binary.Read(bytes.NewReader(data[i*4:]), binary.LittleEndian, &v))
		restored = append(restored, v)
	}

	require.Len(t, restored, len(original))
	for i := range restored {
		assert.Equal(t, original[i], restored[i], "mismatch at index %d", i)
	}
}

func TestFloatsToBytes_512d(t *testing.T) {
	emb := make([]float32, 512)
	for i := range emb {
		emb[i] = float32(i) * 0.001
	}
	data := floatsToBytes(t, emb)
	assert.Len(t, data, 512*4)
}
