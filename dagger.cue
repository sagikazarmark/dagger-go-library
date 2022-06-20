package main

import (
	"dagger.io/dagger"

	"universe.dagger.io/go"

	"github.com/sagikazarmark/dagger-go-library/ci/go/golangci"
)

dagger.#Plan & {
	client: filesystem: ".": read: exclude: [
		".github",
		"bin",
		"build",
		"tmp",
	]
	actions: {
		_source: client.filesystem["."].read.contents

		test: {
			"go": {
				"1.16": _
				"1.17": _
				"1.18": _

				[v=string]: {
					_test: go.#Test & {
						source:  _source
						name:    "go_test_\(v)" // necessary to keep cache for different versions separate
						package: "./..."

						_image: go.#Image & {
							version: v
						}

						input: _image.output
						command: flags: {
							"-race":         true
							"-covermode":    "atomic"
							"-coverprofile": "/coverage.out"
						}

						export: files: "/coverage.out": _
					}
				}
			}
		}
		lint: {
			"go": golangci.#Lint & {
				source:  _source
				version: "1.46"
			}
		}
	}
}
