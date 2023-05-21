//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

const (
	goBin  = "go"
	buf    = "buf"
	goLint = "golangci-lint"
	docker = "docker"
)

// Buf lints and generates code using buf.
func Buf() error {
	if err := sh.Run(buf, "lint"); err != nil {
		return err
	}
	return sh.Run(buf, "generate")
}

// Lint runs golangci-lint.
func Lint() error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, goLint, "run")
	if err != nil {
		return fmt.Errorf("lint failed: %w", err)
	}
	return nil
}

// Test runs the tests.
func Test() error {
	return sh.Run(goBin, "test", "./...")
}

// Docker builds the docker image.
func Docker(service string) error {
	return sh.Run(docker, "build", "-t", service, ".", "--build-arg", "SERVICE="+service)
}

func Docs() error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "mkdocs", "serve")
	if err != nil {
		return fmt.Errorf("docs failed: %w", err)
	}
	return nil
}
