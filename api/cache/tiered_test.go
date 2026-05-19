package cache

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func floatsToBytes(t *testing.T, emb []float32) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	require.NoError(t, binary.Write(buf, binary.LittleEndian, emb))
	return buf.Bytes()
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
