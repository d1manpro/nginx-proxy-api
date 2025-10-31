package nginx

import (
	"fmt"
	"os/exec"
)

func reloadNginx() error {
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx config test failed: %w, output: %s", err, output)
	}

	cmd = exec.Command("nginx", "-s", "reload")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx reload failed: %w, output: %s", err, output)
	}

	return nil
}
