// Package filesystem loads a roadmap tree without following links outside its root.
package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AndySabina/RoadmapControl/internal/domain/roadmap"
)

const maxFileBytes = 1 << 20

type Roadmap struct{ Manifest roadmap.Manifest }

// Load reads roadmap.yaml and its declared YAML modules. Module contents receive
// document-safety checks only; typed module/schema validation belongs to later slices.
//
// This is a cooperative filesystem boundary. It rejects links and unsafe files observed
// during loading, but cannot make check-then-read operations safe against hostile races.
func Load(root string) (Roadmap, error) {
	var result Roadmap
	root, err := cleanRoot(root)
	if err != nil {
		return result, err
	}
	data, err := regularFile(root, filepath.Join(root, "roadmap.yaml"))
	if err != nil {
		return result, fmt.Errorf("roadmap manifest: %w", err)
	}
	manifest, err := roadmap.ParseManifest(data)
	if err != nil {
		return result, err
	}
	declared := make(map[string]bool, len(manifest.Modules))
	for _, name := range manifest.Modules {
		file, err := underRoot(root, name)
		if err != nil {
			return result, err
		}
		data, err := regularFile(root, file)
		if err != nil {
			return result, fmt.Errorf("module %q: %w", name, err)
		}
		if err := roadmap.ValidateYAML(data); err != nil {
			return result, fmt.Errorf("module %q: %w", name, err)
		}
		declared[filepath.ToSlash(name)] = true
	}
	if err := rejectUnlisted(root, declared); err != nil {
		return result, err
	}
	result.Manifest = manifest
	return result, nil
}

func cleanRoot(root string) (string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	info, err := os.Lstat(root)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return "", fmt.Errorf("roadmap root must be a non-symlink directory")
	}
	return root, nil
}

func regularFile(root, name string) ([]byte, error) {
	if err := noSymlinkComponent(root, name); err != nil {
		return nil, err
	}
	info, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("must be a regular file")
	}
	if info.Size() > maxFileBytes {
		return nil, fmt.Errorf("exceeds %d-byte limit", maxFileBytes)
	}
	return os.ReadFile(name)
}

func noSymlinkComponent(root, name string) error {
	rel, err := filepath.Rel(root, name)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return fmt.Errorf("path %q escapes roadmap root", name)
	}
	current := root
	for _, component := range strings.Split(rel, string(filepath.Separator)) {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("must not contain a symlink")
		}
	}
	return nil
}

func underRoot(root, name string) (string, error) {
	file := filepath.Join(root, filepath.FromSlash(name))
	rel, err := filepath.Rel(root, file)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("module path %q escapes roadmap root", name)
	}
	return file, nil
}

func rejectUnlisted(root string, declared map[string]bool) error {
	return filepath.WalkDir(root, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 || entry.IsDir() || !yamlPath(name) || name == filepath.Join(root, "roadmap.yaml") {
			return nil
		}
		rel, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		if !declared[filepath.ToSlash(rel)] {
			return fmt.Errorf("unlisted YAML module %q", filepath.ToSlash(rel))
		}
		return nil
	})
}

func yamlPath(name string) bool {
	return strings.EqualFold(filepath.Ext(name), ".yaml") || strings.EqualFold(filepath.Ext(name), ".yml")
}
