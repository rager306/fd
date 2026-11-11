package queue

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unicode/utf8"
)

func TestResultStoreConcurrentReservation(t *testing.T) {
	s := NewResultStoreWithLimits(4*1024*4, 10)
	t.Cleanup(func() { _ = s.Close() })
	var accepted atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if s.ReservePending(fmt.Sprintf("%016x", i), 1, 1, 1024) {
				accepted.Add(1)
			}
		}(i)
	}
	wg.Wait()
	if accepted.Load() != 4 || s.Size() != 4 {
		t.Fatalf("accepted=%d size=%d, want 4", accepted.Load(), s.Size())
	}
}

func TestResultStoreLifecycle(t *testing.T) {
	s := NewResultStoreWithLimits(8192, 1)
	t.Cleanup(func() { _ = s.Close() })
	id, other := "0000000000000001", "0000000000000002"
	if !s.ReservePending(id, 7, 1, 1024) || s.ReservePending(other, 7, 1, 1024) {
		t.Fatal("entry limit not enforced")
	}
	s.data[id].ts = time.Now().Add(-2 * defaultResultTTL)
	s.evict()
	if _, ok := s.Get(id); !ok {
		t.Fatal("pending reservation expired")
	}
	s.DeletePending(id)
	if !s.ReservePending(other, 7, 1, 1024) {
		t.Fatal("failed enqueue did not release reservation")
	}
	s.CompleteExisting(other, &Result{Status: StatusFailed, Err: errors.New(strings.Repeat("€", 1000))})
	result, ok := s.Get(other)
	if !ok || result.Status != StatusFailed || len(result.Err.Error()) > maxResultErrorBytes || !utf8.ValidString(result.Err.Error()) {
		t.Fatalf("failed result: %+v", result)
	}
	s.DeletePending(other)
	if s.Size() != 1 {
		t.Fatal("deleted terminal result")
	}
	s.data[other].ts = time.Now().Add(-2 * defaultResultTTL)
	if !s.ReservePending(id, 7, 1, 1024) {
		t.Fatal("expired result did not free capacity")
	}
	_ = s.Close()
	if s.ReservePending(other, 7, 1, 1024) {
		t.Fatal("reserved after close")
	}
	s.CompleteExisting(id, &Result{Status: StatusFailed})
	if r, _ := s.Get(id); r.Status != StatusPending {
		t.Fatal("completed after close")
	}
}

func TestResultStoreCopiesAndValidatesResults(t *testing.T) {
	s := NewResultStoreWithLimits(8192, 2)
	t.Cleanup(func() { _ = s.Close() })
	id := "0000000000000001"
	if !s.ReservePending(id, 42, 1, 512) {
		t.Fatal("reserve failed")
	}
	backing := make([]float32, 1024)
	backing[0] = 5
	s.CompleteExisting(id, &Result{Status: StatusCompleted, Dims: 512, Embeddings: [][]float32{backing[:512]}})
	backing[0] = 10
	r, ok := s.Get(id)
	if !ok || r.Status != StatusCompleted || r.Embeddings[0][0] != 5 || cap(r.Embeddings[0]) != 512 || r.CreatedAt != 42 {
		t.Fatalf("result not compact and independent: %+v", r)
	}
	bad := "0000000000000002"
	if !s.ReservePending(bad, 0, 1, 1024) {
		t.Fatal("reserve failed")
	}
	s.CompleteExisting(bad, &Result{Status: StatusCompleted, Dims: 1024, Embeddings: [][]float32{backing[:512]}})
	r, ok = s.Get(bad)
	if !ok || r.Status != StatusFailed || r.Embeddings != nil {
		t.Fatal("accepted malformed result")
	}
	s.CompleteExisting("0000000000000003", &Result{Status: StatusFailed})
	if s.Size() != 2 {
		t.Fatal("completion created unreserved entry")
	}
}
