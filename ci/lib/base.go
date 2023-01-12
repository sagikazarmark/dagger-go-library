package lib

import (
	"crypto/sha256"
	"fmt"

	"dagger.io/dagger"
	"github.com/sagikazarmark/go-option"
)

const (
	defaultGoImageRepository = "docker.io/library/golang"
	defaultGoVersion         = "latest"
)

// Option configures common parameters of all jobs.
type Option interface {
	TestOption
	LintOption
}

// GoVersion specifies which Go version to use.
// GoVersion is ignored when Image is configured.
func GoVersion(v string) Option {
	return goVersion(v)
}

type goVersion string

func (v goVersion) applyTest(o *testOptions) {
	o.baseOptions.GoVersion = string(v)
}

func (v goVersion) applyLint(o *lintOptions) {
	o.baseOptions.GoVersion = string(v)
}

type baseOptions struct {
	ImageRepository string
	GoVersion       string

	Image string

	ProjectRoot string
	CgoEnabled  option.Option[bool]
}

func (o baseOptions) HasCgoValue() bool {
	return o.CgoEnabled != nil && option.IsSome(o.CgoEnabled)
}

func (o baseOptions) CgoValue() string {
	return option.MapOr(o.CgoEnabled, "", func(v bool) string {
		if v {
			return "1"
		}

		return "0"
	})
}

func base(client *dagger.Client, options baseOptions) *dagger.Container {
	projectRoot := "."
	if options.ProjectRoot != "" {
		projectRoot = options.ProjectRoot
	}

	imageRepository := defaultGoImageRepository
	if options.ImageRepository != "" {
		imageRepository = options.ImageRepository
	}

	goVersion := defaultGoVersion
	if options.GoVersion != "" {
		goVersion = options.GoVersion
	}

	image := fmt.Sprintf("%s:%s", imageRepository, goVersion)
	if options.Image != "" {
		image = options.Image
	}

	h := sha256.New()
	h.Write([]byte(image))

	imageHash := h.Sum(nil)

	container := client.Container().
		From(image).
		WithMountedCache("/root/.cache/go-build", client.CacheVolume(fmt.Sprintf("go-build-%x", imageHash))).
		WithMountedCache("/go/pkg/mod", client.CacheVolume(fmt.Sprintf("go-mod-%x", imageHash))).
		WithMountedDirectory("/src", client.Host().Directory(projectRoot)).
		WithWorkdir("/src")

	if options.HasCgoValue() {
		container = container.WithEnvVariable("CGO_ENABLED", options.CgoValue())
	}

	return container
}
