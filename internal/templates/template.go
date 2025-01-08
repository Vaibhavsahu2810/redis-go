package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed html/*
var templateFS embed.FS

type Manager struct {
	templates map[string]*template.Template
}

func New() (*Manager, error) {
	manager := &Manager{
		templates: make(map[string]*template.Template),
	}

	err := fs.WalkDir(templateFS, "html", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		name := strings.TrimSuffix(filepath.Base(path), ".html")

		content, err := templateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading template %s: %w", path, err)
		}

		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("error parsing template %s: %w", path, err)
		}

		manager.templates[name] = tmpl
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error loading templates: %w", err)
	}

	return manager, nil
}

func (m *Manager) Render(name string, data map[string]interface{}) (string, error) {
	tmpl, ok := m.templates[name]
	if !ok {
		return "", fmt.Errorf("template '%s' not found", name)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error rendering template '%s': %w", name, err)
	}

	return buf.String(), nil
}

func (m *Manager) RenderWithSafeURLs(name string, data map[string]interface{}) (string, error) {
	safeData := make(map[string]interface{})
	for key, value := range data {
		safeData[key] = value
	}

	if resetUrl, ok := data["resetUrl"].(string); ok {
		safeData["resetUrl"] = template.URL(resetUrl)
	}

	return m.Render(name, safeData)
}
