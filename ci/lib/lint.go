package lib

import (
	"fmt"

	"dagger.io/dagger"
)

const (
	defaultGolangciLintImageRepository = "docker.io/golangci/golangci-lint"
	defaultGolangciLintVersion         = "latest"
)

type LintOption interface {
	applyLint(o *lintOptions)
}

type lintOptionFunc func(o *lintOptions)

func (f lintOptionFunc) applyLint(o *lintOptions) {
	f(o)
}

type lintOptions struct {
	baseOptions baseOptions

	SourceImageRepository string
	LinterVersion         string

	SourceImage string
}

func Lint(client *dagger.Client, opts ...LintOption) *dagger.Container {
	var options lintOptions

	for _, opt := range opts {
		opt.applyLint(&options)
	}

	sourceImageRepository := defaultGolangciLintImageRepository
	if options.SourceImageRepository != "" {
		sourceImageRepository = options.SourceImageRepository
	}

	linterVersion := defaultGolangciLintVersion
	if options.LinterVersion != "" {
		linterVersion = options.LinterVersion
	}

	sourceImage := fmt.Sprintf("%s:%s", sourceImageRepository, linterVersion)
	if options.SourceImage != "" {
		sourceImage = options.SourceImage
	}

	bin := client.Container().
		From(sourceImage).
		File("/usr/bin/golangci-lint")

	args := []string{"golangci-lint", "run", "--verbose"}

	return base(client, options.baseOptions).
		WithFile("/usr/local/bin/golangci-lint", bin).
		WithExec(args)
}
