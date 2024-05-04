package rsync

import (
	"context"
	"testing"
)

func Test_MultiRSW(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := NewMultiRSW(ctx)

	// Test successful read lock
	if err := r.BeginRead(); err != nil {
		t.Fatalf("Failed to acquire read lock: %v", err)
	}
	r.EndRead()

	// Test successful write lock
	if err := r.BeginWrite(); err != nil {
		t.Fatalf("Failed to acquire write lock: %v", err)
	}
	r.EndWrite()

	// Test that a write blocks other writers and readers.
	err := r.BeginWrite()
	if err != nil {
		t.Fatalf("Failed to acquire write lock in goroutine: %v", err)
	}
	if err := r.BeginRead(); err == nil {
		t.Fatalf("Expected error when reading during active write, got none")
	}
	if err := r.BeginWrite(); err == nil {
		t.Fatalf("Expected error when writing during active write, got none")
	}
	r.EndWrite()

	// Test that a read blocks a writer.
	err = r.BeginRead()
	if err != nil {
		t.Fatalf("Failed to acquire read lock in goroutine: %v", err)
	}
	if err := r.BeginWrite(); err == nil {
		t.Fatalf("Expected error when writing during active read, got none")
	}
	r.EndRead()

	// Test that a reader doesn't block other readers.
	err = r.BeginRead()
	if err != nil {
		t.Fatalf("Failed to acquire read lock in goroutine: %v", err)
	}
	if err := r.BeginRead(); err != nil {
		t.Fatalf("Failed to acquire read lock in goroutine: %v", err)
	}
	r.EndRead()
	r.EndRead()
}
