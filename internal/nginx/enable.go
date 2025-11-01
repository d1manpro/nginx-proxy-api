package nginx

import (
	"fmt"
	"os"
	"path/filepath"
)

func activateSite(fileName string) error {
	src := filepath.Join(sitesAvailable, fileName)
	dst := filepath.Join(sitesEnabled, fileName)

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("config-file %s is not found", src)
	}

	if _, err := os.Lstat(dst); os.IsNotExist(err) {
		if err := os.Symlink(src, dst); err != nil {
			return fmt.Errorf("failed to create symlink: %v", err)
		}
	} else {
		fmt.Println("symlink is already exists")
	}
	return reloadNginx()
}
