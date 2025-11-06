package certbot

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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

func IsCertExists(domain string) (bool, error) {
	cmd := exec.Command( "certbot", "certificates")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("certbot get-cert-list error: %v\n%s", err, out.String())
	}

	return strings.Contains(out.String(), domain), nil
}
