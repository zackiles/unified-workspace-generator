//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Build binaries for all platforms
func Build() error {
	platforms := []string{"linux/amd64", "darwin/amd64", "windows/amd64"}

	for _, platform := range platforms {
		fmt.Println("Building for", platform)
		err := buildForPlatform(platform)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildForPlatform(platform string) error {
	var outputFile string
	parts := strings.Split(platform, "/")
	GOOS := parts[0]
	GOARCH := parts[1]

	if GOOS == "windows" {
		outputFile = "create-project.exe"
	} else {
		outputFile = "create-project"
	}

	cmd := exec.Command("go", "build", "-o", "bin/"+GOOS+"/"+outputFile)
	cmd.Env = append(os.Environ(), "GOOS="+GOOS, "GOARCH="+GOARCH)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
