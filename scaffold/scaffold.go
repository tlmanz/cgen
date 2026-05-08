package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const originalModule = "github.com/tlmanz/catalyst/v3"

// New scaffolds a new Catalyst project with the given Go module name.
func New(module, outDir string, skipConfirm bool) error {
	parts := strings.Split(strings.TrimRight(module, "/"), "/")
	projectName := parts[len(parts)-1]
	projectDir := filepath.Join(outDir, projectName)

	if !skipConfirm {
		fmt.Printf("Creating project '%s' at '%s'\n", projectName, projectDir)
		fmt.Print("Continue? [y/N]: ")
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return fmt.Errorf("creating project directory: %w", err)
	}

	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel("templates", path)
		if rel == "." {
			return nil
		}

		dest := filepath.Join(projectDir, rel)

		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}

		// strip .tmpl suffix so go.mod.tmpl → go.mod, main.go.tmpl → main.go, etc.
		dest = strings.TrimSuffix(dest, ".tmpl")

		return writeFile(path, dest, module, d)
	})
	if err != nil {
		return err
	}

	// write the real project name so metadata.Name() returns the correct value
	// from the very first build, before set_metadata.sh is run.
	namePath := filepath.Join(projectDir, "metadata", "name.txt")
	if err := os.WriteFile(namePath, []byte(projectName), 0o644); err != nil {
		return fmt.Errorf("writing metadata/name.txt: %w", err)
	}

	fmt.Printf("Project '%s' created successfully at '%s'\n", projectName, projectDir)
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectDir)
	fmt.Println("  go mod tidy")
	fmt.Println("  cp config.example.yaml config.yaml")
	fmt.Println("  go run ./cmd/server")
	return nil
}

func writeFile(src, dest, module string, _ fs.DirEntry) error {
	data, err := templateFS.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading template %s: %w", src, err)
	}

	if isText(src) {
		data = []byte(strings.ReplaceAll(string(data), originalModule, module))
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	return os.WriteFile(dest, data, 0o644)
}

func isText(path string) bool {
	textSuffixes := []string{
		".go", ".mod", ".sum", ".tmpl",
		".yaml", ".yml", ".toml", ".json",
		".sh", ".md", ".txt", ".env",
		"Makefile", "Dockerfile",
		".gitignore", ".dockerignore", ".editorconfig",
	}
	for _, s := range textSuffixes {
		if strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
}
