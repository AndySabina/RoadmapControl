package filesystem

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const testManifest = "schema: https://schemas.example/v1\nmodules: [%s]\n"

func TestLoadReadsNestedDeclaredModulesWithoutWriting(t *testing.T) {
	root := t.TempDir()
	write(t, root, "roadmap.yaml", manifest("nested/one.yaml, two.YML"))
	write(t, root, "nested/one.yaml", "name: one\n")
	write(t, root, "two.YML", "name: two\n")

	got, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Manifest.Modules) != 2 {
		t.Fatalf("loaded modules = %#v", got.Manifest.Modules)
	}
	for name, want := range map[string]string{
		"nested/one.yaml": "name: one\n",
		"two.YML":         "name: two\n",
	} {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || string(data) != want {
			t.Fatalf("module %q changed: got %q, err %v", name, data, err)
		}
	}
}

func TestLoadRejectsUnsafeFilesystemInputs(t *testing.T) {
	cases := []struct {
		name  string
		setup func(t *testing.T, root string) string
	}{
		{"missing module", func(t *testing.T, root string) string {
			write(t, root, "roadmap.yaml", manifest("missing.yaml"))
			return root
		}},
		{"unlisted lowercase YAML", func(t *testing.T, root string) string {
			write(t, root, "roadmap.yaml", manifest("listed.yaml"))
			write(t, root, "listed.yaml", "ok: true\n")
			write(t, root, "extra.yaml", "ok: true\n")
			return root
		}},
		{"unlisted uppercase YML", func(t *testing.T, root string) string {
			write(t, root, "roadmap.yaml", manifest("listed.yaml"))
			write(t, root, "listed.yaml", "ok: true\n")
			write(t, root, "extra.YML", "ok: true\n")
			return root
		}},
		{"manifest directory", func(t *testing.T, root string) string {
			mkdir(t, filepath.Join(root, "roadmap.yaml"))
			return root
		}},
		{"declared module directory", func(t *testing.T, root string) string {
			write(t, root, "roadmap.yaml", manifest("module.yaml"))
			mkdir(t, filepath.Join(root, "module.yaml"))
			return root
		}},
		{"root symlink", func(t *testing.T, root string) string {
			outside := t.TempDir()
			write(t, outside, "roadmap.yaml", manifest(""))
			return symlink(t, outside, filepath.Join(root, "linked-root"))
		}},
		{"manifest symlink", func(t *testing.T, root string) string {
			outside := filepath.Join(t.TempDir(), "roadmap.yaml")
			mustWrite(t, outside, manifest(""))
			symlink(t, outside, filepath.Join(root, "roadmap.yaml"))
			return root
		}},
		{"declared module symlink", func(t *testing.T, root string) string {
			write(t, root, "roadmap.yaml", manifest("module.yaml"))
			outside := filepath.Join(t.TempDir(), "module.yaml")
			mustWrite(t, outside, "ok: true\n")
			symlink(t, outside, filepath.Join(root, "module.yaml"))
			return root
		}},
		{"intermediate directory symlink", func(t *testing.T, root string) string {
			write(t, root, "roadmap.yaml", manifest("nested/module.yaml"))
			outside := t.TempDir()
			write(t, outside, "module.yaml", "ok: true\n")
			symlink(t, outside, filepath.Join(root, "nested"))
			return root
		}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Load(tt.setup(t, t.TempDir())); err == nil {
				t.Fatal("Load accepted unsafe filesystem input")
			}
		})
	}
}

func TestLoadRejectsUnsafeOrOversizedYAML(t *testing.T) {
	for _, tt := range []struct {
		name, module string
	}{
		{"second document", "one: 1\n---\ntwo: 2\n"},
		{"alias", "one: &one value\ntwo: *one\n"},
		{"duplicate key", "one: 1\none: 2\n"},
		{"non UTF-8", "one: \xff\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			write(t, root, "roadmap.yaml", manifest("module.yaml"))
			write(t, root, "module.yaml", tt.module)
			if _, err := Load(root); err == nil {
				t.Fatal("Load accepted unsafe YAML")
			}
		})
	}
	for _, tt := range []struct{ name, path string }{{"manifest", "roadmap.yaml"}, {"module", "module.yaml"}} {
		t.Run("oversized "+tt.name, func(t *testing.T) {
			root := t.TempDir()
			if tt.path == "roadmap.yaml" {
				write(t, root, tt.path, string(bytes.Repeat([]byte("x"), 1<<20+1)))
			} else {
				write(t, root, "roadmap.yaml", manifest("module.yaml"))
				write(t, root, tt.path, string(bytes.Repeat([]byte("x"), 1<<20+1)))
			}
			if _, err := Load(root); err == nil {
				t.Fatal("Load accepted oversized YAML")
			}
		})
	}
}

func manifest(modules string) string { return fmt.Sprintf(testManifest, modules) }

func write(t *testing.T, root, name, body string) {
	t.Helper()
	path := filepath.Join(root, name)
	mkdir(t, filepath.Dir(path))
	mustWrite(t, path, body)
}
func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Fatal(err)
	}
}
func symlink(t *testing.T, target, link string) string {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	return link
}
