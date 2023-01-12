package lib

import (
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

type testOptions struct {
	baseOptions

	Verbose             bool
	RaceDetectorEnabled bool
}

func Test(client *dagger.Client, opts ...TestOption) *dagger.Container {
	var options testOptions

	for _, opt := range opts {
		opt.applyTest(&options)
	}

	args := []string{"go", "test"}

	if options.Verbose {
		args = append(args, "-v")
	}

	if options.RaceDetectorEnabled {
		args = append(args, "-race")

		options.baseOptions.CgoEnabled = option.Some(true)
	}

	args = append(args, "./...")

	return base(client, options.baseOptions).WithExec(args)
}
