package helm

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Chart2Deploy(chartFile string) (string, error) {

	command := fmt.Sprint("helm install --dry-run tmp %s", chartFile)

	// helm
	cmd := exec.Command(command)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	if stderr.Len() > 0 {
		return stderr.String(), nil
	}

	return stdout.String(), nil

}
