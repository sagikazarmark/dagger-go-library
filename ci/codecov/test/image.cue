package codecov

import (
	"dagger.io/dagger"

	"github.com/sagikazarmark/dagger-go-library/ci/codecov"
)

dagger.#Plan & {
	actions: test: codecov.#Image & {}
}
