package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsTempFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{".hidden", true},
		{"file~", true},
		{"file.swp", true},
		{".file.swp", true},
		{"normal.go", false},
		{"README.md", false},
		{"main_test.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTempFile(tt.name); got != tt.want {
				t.Errorf("isTempFile(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsSkippedDir(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/project/.git/objects", true},
		{"/project/node_modules/pkg", true},
		{"/project/vendor/lib", true},
		{"/project/src/main.go", false},
		{"/project/docs", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isSkippedDir(tt.path); got != tt.want {
				t.Errorf("isSkippedDir(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestResolveWatchDirs(t *testing.T) {
	base := "/tmp/watchdoc-test-base"

	t.Run("empty extras", func(t *testing.T) {
		watchList, extrasList := resolveWatchDirs("", base)
		if len(watchList) != 1 || watchList[0] != base {
			t.Errorf("expected [%s], got %v", base, watchList)
		}
		if extrasList != nil {
			t.Errorf("expected nil extras, got %v", extrasList)
		}
	})

	t.Run("with extras", func(t *testing.T) {
		dir1, err := os.MkdirTemp("", "wd1")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.RemoveAll(dir1) }()
		dir2, err := os.MkdirTemp("", "wd2")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.RemoveAll(dir2) }()

		watchList, extrasList := resolveWatchDirs(dir1+","+dir2, base)
		if len(watchList) != 3 {
			t.Errorf("expected 3 watch dirs, got %d: %v", len(watchList), watchList)
		}
		if len(extrasList) != 2 {
			t.Errorf("expected 2 extras, got %d: %v", len(extrasList), extrasList)
		}
	})

	t.Run("whitespace handling", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "wd")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.RemoveAll(dir) }()

		watchList, _ := resolveWatchDirs("  "+dir+" , , ", base)
		if len(watchList) != 2 {
			t.Errorf("expected 2 (base + 1 valid), got %d: %v", len(watchList), watchList)
		}
	})
}

func TestIsSourceFile(t *testing.T) {
	sourceDirs := []string{"/project/src", "/project/lib"}

	tests := []struct {
		path string
		want bool
	}{
		{"/project/src/main.go", true},
		{"/project/lib/util.go", true},
		{"/project/build/out.html", false},
		{"/other/src/file.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isSourceFile(tt.path, sourceDirs); got != tt.want {
				t.Errorf("isSourceFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsOutputFile(t *testing.T) {
	servedDir := "/project/build"

	tests := []struct {
		path string
		want bool
	}{
		{"/project/build/index.html", true},
		{"/project/build/sub/page.html", true},
		{"/project/src/main.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isOutputFile(tt.path, servedDir); got != tt.want {
				t.Errorf("isOutputFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsTempFile_WithPath(t *testing.T) {
	// Test with full paths to ensure filepath.Base is used
	tests := []struct {
		path string
		want bool
	}{
		{filepath.Join("/some/dir", ".hidden"), true},
		{filepath.Join("/some/dir", "file.txt~"), true},
		{filepath.Join("/some/dir", "normal.go"), false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isTempFile(tt.path); got != tt.want {
				t.Errorf("isTempFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
