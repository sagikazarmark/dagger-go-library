//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"

	"github.com/magefile/mage/mg"
	"github.com/sagikazarmark/dagger-go-library/ci/lib"
)

// Run tests
func Test(ctx context.Context) error {
	var clientOpts []dagger.ClientOpt

	if os.Getenv("DEBUG") == "true" {
		clientOpts = append(clientOpts, dagger.WithLogOutput(os.Stderr))
	}

	client, err := dagger.Connect(ctx, clientOpts...)
	if err != nil {
		return err
	}
	defer client.Close()

	var opts []lib.TestOption

	if goVersion := os.Getenv("GO_VERSION"); goVersion != "" {
		opts = append(opts, lib.GoVersion(goVersion))
	}

	return process(ctx, lib.Test(client, opts...))
}

// Run linter
func Lint(ctx context.Context) error {
	var clientOpts []dagger.ClientOpt

	if os.Getenv("DEBUG") == "true" {
		clientOpts = append(clientOpts, dagger.WithLogOutput(os.Stderr))
	}

	client, err := dagger.Connect(ctx, clientOpts...)
	if err != nil {
		return err
	}
	defer client.Close()

	var opts []lib.LintOption

	if goVersion := os.Getenv("GO_VERSION"); goVersion != "" {
		opts = append(opts, lib.GoVersion(goVersion))
	}

	return process(ctx, lib.Lint(client, opts...))
}

func process(ctx context.Context, container *dagger.Container) error {
	output, err := container.Stdout(ctx)

	fmt.Print(output)

	// if err != nil {
	// 	return err
	// }

	erroutput, err := container.Stderr(ctx)

	fmt.Print(erroutput)

	if err != nil {
		return err
	}

	exit, err := container.ExitCode(ctx)
	if err != nil {
		return err
	}

	if exit > 0 {
		return mg.Fatal(exit)
	}

	return nil
}
