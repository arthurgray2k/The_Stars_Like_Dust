package comm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arthurgray2k/The_Stars_Like_Dust/pkg/dsl"
)

func TestCommManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "comm_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr := NewManager(tempDir)

	// Initially empty
	items, err := mgr.LoadAll()
	if err != nil {
		t.Fatalf("unexpected error loading empty dir: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}

	// Save a dispatch
	disp := &dsl.Dispatch{
		ID:             "HCP-TEST-001",
		Encryption:     "NONE",
		Origin:         "Earth",
		Target:         "Biron-Farrill",
		Stardate:       "820-G.E.",
		Priority:       "URGENT",
		Payload:        []dsl.PayloadItem{{Type: "INTEL", Content: "Test alert"}},
		Authentication: "SIG-TEST",
	}

	err = mgr.SaveDispatch("test_001.hcp", disp)
	if err != nil {
		t.Fatalf("failed to save dispatch: %v", err)
	}

	// Save invalid file to test parse error
	badFile := filepath.Join(tempDir, "invalid.hcp")
	_ = os.WriteFile(badFile, []byte("NOT A VALID HCP"), 0644)

	_, err = mgr.LoadAll()
	if err == nil {
		t.Errorf("expected parse error for invalid.hcp, got nil")
	}
	_ = os.Remove(badFile)

	// Also non-hcp file should be ignored
	nonHcp := filepath.Join(tempDir, "notes.txt")
	_ = os.WriteFile(nonHcp, []byte("some notes"), 0644)

	// Read back
	items, err = mgr.LoadAll()
	if err != nil {
		t.Fatalf("failed to reload dispatches: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 dispatch, got %d", len(items))
	}
	if items[0].ID != "HCP-TEST-001" {
		t.Errorf("expected ID HCP-TEST-001, got %s", items[0].ID)
	}

	// Test FilterByTarget
	filtered := FilterByTarget(items, "Biron")
	if len(filtered) != 1 {
		t.Errorf("expected 1 match for Biron, got %d", len(filtered))
	}

	filteredBroadcast := FilterByTarget([]*dsl.Dispatch{
		{Target: "BROADCAST"},
		{Target: "ALL"},
	}, "AnyTarget")
	if len(filteredBroadcast) != 2 {
		t.Errorf("expected 2 matches for broadcast/all, got %d", len(filteredBroadcast))
	}

	filteredNone := FilterByTarget(items, "Aratap")
	if len(filteredNone) != 0 {
		t.Errorf("expected 0 matches for Aratap, got %d", len(filteredNone))
	}

	// Test non-existent directory
	nonExistentMgr := NewManager(filepath.Join(tempDir, "non_existent"))
	items, err = nonExistentMgr.LoadAll()
	if err != nil || len(items) != 0 {
		t.Errorf("expected empty results without error for non-existent dir, got %v, %d items", err, len(items))
	}
}

func TestCommManagerExistingRepoDispatches(t *testing.T) {
	commDir := filepath.Join("..", "..", "comm")
	mgr := NewManager(commDir)
	dispatches, err := mgr.LoadAll()
	if err != nil {
		t.Fatalf("failed to load real project dispatches: %v", err)
	}
	if len(dispatches) < 5 {
		t.Errorf("expected at least 5 dispatches in comm/, got %d", len(dispatches))
	}
}
