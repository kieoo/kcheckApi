package util

import (
	"bufio"
	"bytes"
	"path/filepath"
	"regexp"
)

func CleanOriYaml(ori []byte) map[string][]byte {
	// 只保留deployment, statefulset的配置, 并切片
	SIndex := bytes.Index(ori, []byte("# Source"))
	if SIndex == -1 {
		return nil
	}
	deploy := ori[SIndex:]

	out := make(map[string][]byte)

	for i, v := range bytes.Split(deploy, []byte("# Source: ")) {
		if i == 0 {
			continue
		}
		matchDeploy, _ := regexp.Match(`kind:\s*Deployment`, v)
		matchState, _ := regexp.Match(`kind:\s*StatefulSet`, v)
		if !matchDeploy && !matchState {
			continue
		}
		bytesReader := bytes.NewReader(v)
		bufReader := bufio.NewReader(bytesReader)
		filePath, _, _ := bufReader.ReadLine()
		filePath = bytes.TrimSpace(filePath)
		fileName := filepath.Base(string(filePath))
		content := append([]byte("# Source: "), v[:]...)
		out[fileName] = content
	}

	return out
}
