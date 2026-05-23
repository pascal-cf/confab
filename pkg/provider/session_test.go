package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathIsUnderAnyRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "root")
	other := filepath.Join(tmpDir, "other")
	for _, d := range []string{root, other} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("MkdirAll(%s): %v", d, err)
		}
	}

	t.Run("under single root", func(t *testing.T) {
		path := filepath.Join(root, "sub", "file.jsonl")
		if !pathIsUnderAnyRoot(path, []string{root}) {
			t.Error("expected true for path under root")
		}
	})

	t.Run("not under root", func(t *testing.T) {
		path := filepath.Join(other, "file.jsonl")
		if pathIsUnderAnyRoot(path, []string{root}) {
			t.Error("expected false for path outside root")
		}
	})

	t.Run("multiple roots — matches second", func(t *testing.T) {
		path := filepath.Join(other, "file.jsonl")
		if !pathIsUnderAnyRoot(path, []string{root, other}) {
			t.Error("expected true when path matches second root")
		}
	})

	t.Run("traversal via .. rejected by caller, cleaned path is safe", func(t *testing.T) {
		// The caller (ValidateTranscriptPath) rejects ".." before calling us.
		// Verify that a lexically-cleaned path that stays under root is accepted.
		path := filepath.Clean(filepath.Join(root, "sub", "..", "file.jsonl"))
		if !pathIsUnderAnyRoot(path, []string{root}) {
			t.Error("expected true for cleaned path under root")
		}
	})

	t.Run("nonexistent parent — lexical fallback", func(t *testing.T) {
		// Parent directory does not exist; EvalSymlinks fails, falls back to lexical check.
		path := filepath.Join(root, "newproject", "session.jsonl")
		if !pathIsUnderAnyRoot(path, []string{root}) {
			t.Error("expected true for nonexistent parent under root (lexical fallback)")
		}
	})

	t.Run("symlinked parent escaping root", func(t *testing.T) {
		linkDir := filepath.Join(root, "link-out")
		if err := os.Symlink(other, linkDir); err != nil {
			t.Fatalf("Symlink: %v", err)
		}
		path := filepath.Join(linkDir, "file.jsonl")
		if pathIsUnderAnyRoot(path, []string{root}) {
			t.Error("expected false for symlink escaping root")
		}
	})
}
