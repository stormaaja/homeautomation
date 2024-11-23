package store

import (
	"testing"
)

func TestMemoryStore(t *testing.T) {
	store := &MemoryStore{Data: make(map[string]interface{})}

	// Test SetString and GetString
	store.SetString("key1", "value1")
	if val, _ := store.GetString("key1"); val != "value1" {
		t.Errorf("expected 'value1', got '%s'", val)
	}

	// Test SetInt and GetInt
	store.SetInt("key2", 42)
	if val, _ := store.GetInt("key2"); val != 42 {
		t.Errorf("expected 42, got %d", val)
	}

	// Test SetFloat and GetFloat
	store.SetFloat("key3", 3.14)
	if val, _ := store.GetFloat("key3"); val != 3.14 {
		t.Errorf("expected 3.14, got %f", val)
	}
}

func TestMemoryStoreGetNonExistingKey(t *testing.T) {
	store := &MemoryStore{Data: make(map[string]interface{})}

	// Test GetString
	if val, ok := store.GetString("key"); ok {
		t.Errorf("expected ok to be false, got true with value '%s'", val)
	}

	// Test GetInt
	if val, ok := store.GetInt("key"); ok {
		t.Errorf("expected ok to be false, got true with value %d", val)
	}

	// Test GetFloat
	if val, ok := store.GetFloat("key"); ok {
		t.Errorf("expected ok to be false, got true with value %f", val)
	}
}
