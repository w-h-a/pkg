package kubernetes

import (
	"os"
	"path"
)

func DetectNamespace() (string, error) {
	nsPath := path.Join(serviceAccountPath, "namespace")

	// make sure it's a file and that we can read it
	if file, err := os.Stat(nsPath); err != nil {
		return "", err
	} else if file.IsDir() {
		return "", ErrReadNamespace
	}

	// read it
	ns, err := os.ReadFile(nsPath)
	if err != nil {
		return "", err
	}

	return string(ns), nil
}