package cache

import (
	"context"
	"testing"
	"time"
)

func TestRedisClient_NewRedisClient(t *testing.T) {
	// This test is skipped in CI/CD as it requires a running Redis instance
	if testing.Short() {
		t.Skip("Skipping Redis integration test in short mode")
	}

	client, err := NewRedisClient("localhost", "6379", "", 10)
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
		return
	}
	defer client.Close()

	if client == nil {
		t.Fatal("Expected client to be created")
	}
}

func TestRedisClient_SetAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis integration test in short mode")
	}

	client, err := NewRedisClient("localhost", "6379", "", 10)
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	// Test Set and Get
	key := "test:key:123"
	value := "test-value"

	err = client.Set(ctx, key, value)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	retrieved, err := client.Get(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	if retrieved != value {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}

	// Cleanup
	_ = client.Del(ctx, key)
}

func TestRedisClient_Del(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis integration test in short mode")
	}

	client, err := NewRedisClient("localhost", "6379", "", 10)
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	// Set a value
	key := "test:del:123"
	_ = client.Set(ctx, key, "value")

	// Delete it
	err = client.Del(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// Verify it's gone
	_, err = client.Get(ctx, key)
	if err == nil {
		t.Error("Expected error for missing key, got nil")
	}
}

func TestRedisClient_DeleteByPattern(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis integration test in short mode")
	}

	client, err := NewRedisClient("localhost", "6379", "", 10)
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	// Set multiple keys with same pattern
	keys := []string{
		"test:pattern:1",
		"test:pattern:2",
		"test:pattern:3",
	}

	for _, key := range keys {
		_ = client.Set(ctx, key, "value")
	}

	// Delete by pattern
	err = client.DeleteByPattern(ctx, "test:pattern:*")
	if err != nil {
		t.Fatalf("Failed to delete by pattern: %v", err)
	}

	// Give Redis a moment to process
	time.Sleep(100 * time.Millisecond)

	// Verify all are gone
	for _, key := range keys {
		_, err := client.Get(ctx, key)
		if err == nil {
			t.Errorf("Expected key %s to be deleted", key)
		}
	}
}

func TestRedisClient_TTL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis integration test in short mode")
	}

	// Use short TTL for testing
	client, err := NewRedisClient("localhost", "6379", "", 0) // 0 minutes = use seconds
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
		return
	}

	// Override TTL to 1 second for testing
	client.ttl = 1 * time.Second
	defer client.Close()

	ctx := context.Background()

	key := "test:ttl:123"
	_ = client.Set(ctx, key, "value")

	// Verify it exists
	_, err = client.Get(ctx, key)
	if err != nil {
		t.Fatal("Expected key to exist")
	}

	// Wait for TTL to expire
	time.Sleep(2 * time.Second)

	// Verify it's gone
	_, err = client.Get(ctx, key)
	if err == nil {
		t.Error("Expected key to be expired")
	}
}
