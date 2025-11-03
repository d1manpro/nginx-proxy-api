package nginx

import (
	"fmt"
	"html/template"
	"os"
)

func AddConfig(domain, certDomain, target, tmplStr, fileName string) error {
	tmpl, err := template.New("nginx").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create("/etc/nginx/sites-available/" + fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	cfg := tmplConfig{
		Domain:     domain,
		CertDomain: certDomain,
		Target:     target,
	}

	if err := tmpl.Execute(file, cfg); err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	return activateSite(fileName)
}

func RemoveConfig(fileName string) error {
	paths := []string{
		"/etc/nginx/sites-enabled/" + fileName,
		"/etc/nginx/sites-available/" + fileName,
	}

	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	return reloadNginx()
}
