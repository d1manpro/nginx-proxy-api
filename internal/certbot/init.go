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
		return fmt.Errorf("certbot error: %w, output: %s", err, output)
	}
	return nil
}
