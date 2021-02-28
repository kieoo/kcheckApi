package helm

import (
	"bytes"
	"os/exec"
)

func Chart2Deploy(chartFile string) ([]byte, error) {

	// helm

	cmd := exec.Command("helm", "install", "--dry-run", "tmp", chartFile)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stderr.Bytes(), err
	}

	return Formatter(stdout.Bytes())

}

func Formatter(chartOutPut []byte) ([]byte, error) {

	SIndex := bytes.Index(chartOutPut, []byte("# Source"))
	deploy := chartOutPut[SIndex:]
	return deploy, nil
}
