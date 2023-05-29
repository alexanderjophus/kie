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

var (
	services = []string{"ingestion", "features", "cleaning"}
)

// Buf lints and generates code using buf.
func Buf() error {
	if err := run(buf, "lint"); err != nil {
		return err
	}
	return run(buf, "generate")
}

// Lint runs golangci-lint.
func Lint() error {
	return run(goLint, "run")
}

// Test runs the tests.
func Test() error {
	return run(goBin, "test", "./...")
}

// Docker builds the docker image.
func Docker(service string) error {
	return run(docker, "build", "-t", service, ".", "--build-arg", "SERVICE="+service)
}

func ensurePachDRunning() error {
	if err := run("helm", "repo", "add", "pachyderm", "https://pachyderm.github.io/helm-chart"); err != nil {
		return err
	}
	if err := run("helm", "repo", "update"); err != nil {
		return err
	}
	if err := run("helm", "install", "pachd", "pachyderm/pachyderm", "--create-namespace", "-n", "pachd", "--set", "deployTarget=LOCAL", "--set", "proxy.enabled=false"); err != nil {
		return err
	}
	return nil
}

func Deploy() error {
	// if err := ensurePachDRunning(); err != nil {
	// 	return err
	// }
	return run("kubectl", "apply", "-k", "deploy/base")
}

func BuildDeploy() error {
	for _, service := range services {
		if err := Docker(service); err != nil {
			return err
		}
	}
	return Deploy()
}

// Docs serves the docs locally.
func Docs() error {
	return run("mkdocs", "serve")
}

func run(cmd string, args ...string) error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, cmd, args...)
	if err != nil {
		return fmt.Errorf("%s failed: %w", cmd, err)
	}
	return nil
}
