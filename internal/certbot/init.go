package certbot

import (
	"fmt"
	"os/exec"
)

func GetCert(domain, email string) error {
	cmd := exec.Command("certbot", "certonly",
		"--nginx",
		"-d", domain,
		"--agree-tos",
		"--non-interactive",
		"-m", email,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("certbot get error: %w, output: %s", err, output)
	}
	return nil
}

func DeleteCert(domain string) error {
	cmd := exec.Command("certbot", "delete",
		"--cert-name", domain,
		"--non-interactive",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("certbot delete error: %w, output: %s", err, output)
	}

	return nil
}
