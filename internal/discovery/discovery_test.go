package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRootFromExactDir(t *testing.T) {
	dir := t.TempDir()
	roadmapDir := filepath.Join(dir, ".liste")
	os.MkdirAll(roadmapDir, 0755)

	result, err := FindRoot(dir)
	if err != nil {
		t.Fatalf("FindRoot failed: %v", err)
	}
	if result != roadmapDir {
		t.Errorf("FindRoot = %q, want %q", result, roadmapDir)
	}
}

func TestFindRootFromSubdir(t *testing.T) {
	dir := t.TempDir()
	roadmapDir := filepath.Join(dir, ".liste")
	os.MkdirAll(roadmapDir, 0755)

	subDir := filepath.Join(dir, "sub", "deep")
	os.MkdirAll(subDir, 0755)

	result, err := FindRoot(subDir)
	if err != nil {
		t.Fatalf("FindRoot failed: %v", err)
	}
	if result != roadmapDir {
		t.Errorf("FindRoot = %q, want %q", result, roadmapDir)
	}
}

func TestFindRootNotFound(t *testing.T) {
	dir := t.TempDir()

	result, err := FindRoot(dir)
	if err != nil {
		t.Fatalf("FindRoot failed: %v", err)
	}
	if result != "" {
		t.Errorf("FindRoot = %q, want empty string", result)
	}
}

func TestFindSubProjects(t *testing.T) {
	dir := t.TempDir()

	// Create root .liste
	os.MkdirAll(filepath.Join(dir, ".liste"), 0755)

	// Create sub-project .liste directories
	os.MkdirAll(filepath.Join(dir, "service-a", ".liste"), 0755)
	os.MkdirAll(filepath.Join(dir, "service-b", ".liste"), 0755)
	os.MkdirAll(filepath.Join(dir, "nested", "service-c", ".liste"), 0755)

	subs, err := FindSubProjects(dir)
	if err != nil {
		t.Fatalf("FindSubProjects failed: %v", err)
	}

	if len(subs) != 3 {
		t.Fatalf("Found %d sub-projects, want 3", len(subs))
	}

	names := make(map[string]bool)
	for _, sub := range subs {
		names[sub.Name] = true
	}

	if !names["service-a"] {
		t.Error("Missing sub-project: service-a")
	}
	if !names["service-b"] {
		t.Error("Missing sub-project: service-b")
	}
	if !names["nested/service-c"] {
		t.Error("Missing sub-project: nested/service-c")
	}
}

func TestFindSubProjectsSkipsHiddenDirs(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, ".liste"), 0755)
	os.MkdirAll(filepath.Join(dir, ".hidden", ".liste"), 0755)

	subs, err := FindSubProjects(dir)
	if err != nil {
		t.Fatalf("FindSubProjects failed: %v", err)
	}

	if len(subs) != 0 {
		t.Errorf("Found %d sub-projects, want 0 (should skip .hidden)", len(subs))
	}
}

func TestDiscoverFull(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, ".liste"), 0755)
	os.MkdirAll(filepath.Join(dir, "svc", ".liste"), 0755)

	result, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}
	if result == nil {
		t.Fatal("Discover returned nil")
	}
	if result.Root != filepath.Join(dir, ".liste") {
		t.Errorf("Root = %q, want %q", result.Root, filepath.Join(dir, ".liste"))
	}
	if len(result.SubProjects) != 1 {
		t.Errorf("SubProjects = %d, want 1", len(result.SubProjects))
	}
}

func TestDiscoverNone(t *testing.T) {
	dir := t.TempDir()

	result, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}
	if result != nil {
		t.Errorf("Discover = %+v, want nil", result)
	}
}
