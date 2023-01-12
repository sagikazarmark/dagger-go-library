package lib

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/sagikazarmark/go-option"
)

type TestOption interface {
	applyTest(o *testOptions)
}

type testOptionFunc func(o *testOptions)

func (f testOptionFunc) applyTest(o *testOptions) {
	f(o)
}

func Verbose(v bool) TestOption {
	return testOptionFunc(func(o *testOptions) {
		o.Verbose = v
	})
}

func EnableRaceDetector() TestOption {
	return testOptionFunc(func(o *testOptions) {
		o.RaceDetectorEnabled = true
	})
}

func EnableCoverage() TestOption {
	return testOptionFunc(func(o *testOptions) {
		o.CoverageEnabled = true
	})
}

type testOptions struct {
	baseOptions

	Verbose             bool
	RaceDetectorEnabled bool
	CoverageEnabled     bool
}

func Test(client *dagger.Client, opts ...TestOption) *dagger.Container {
	var options testOptions

	for _, opt := range opts {
		opt.applyTest(&options)
	}

	var args []string

	if false {
		args = []string{"gotestsum", "--no-summary=skipped", "--junitfile", "coverage.json"}
	} else {
		args = []string{"go", "test"}
	}

	if options.Verbose {
		args = append(args, "-v")
	}

	if options.RaceDetectorEnabled {
		args = append(args, "-race")

		options.baseOptions.CgoEnabled = option.Some(true)
	}

	if options.CoverageEnabled {
		args = append(args, "-coverprofile=coverage.txt", "-covermode=atomic")
	}

	args = append(args, "./...")

	return base(client, options.baseOptions).WithExec(args)
}

func downloadGotestsum(container *dagger.Container, version string) *dagger.Container {
	url := fmt.Sprintf("https://github.com/gotestyourself/gotestsum/releases/download/v%s/gotestsum_%s_linux_amd64.tar.gz", version, version)

	return container.
		WithExec([]string{"sh", "-c", fmt.Sprintf("wget -O gotestsum.tar.gz %s", url)}).
		// WithExec([]string{"sh", "-c", "echo \"$MEMCACHED_SHA1  memcached.tar.gz\" | sha1sum -c -"}).
		WithExec([]string{"mkdir", "-p", "/usr/src/gotestsum"}).
		WithExec([]string{"tar", "-xvf", "gotestsum.tar.gz", "-C", "/usr/src/gotestsum", "--strip-components=1"}).
		WithExec([]string{"rm", "gotestsum.tar.gz"})
}
