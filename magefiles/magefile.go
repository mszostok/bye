package main

import (
	"github.com/magefile/mage/mg"
	"github.com/samber/lo"
	"go.szostok.io/magex/deps"
	"go.szostok.io/magex/shx"
)

const (
	golangciLintVersion = "1.49.0"
	bin                 = "bin"
)

var (
	Aliases = map[string]interface{}{
		"l": Lint,
	}
)

// "Go" Targets

// Lint Runs linters on the codebase
func Lint() error {
	lo.Must0(deps.EnsureGolangciLint(bin, golangciLintVersion))
	return shx.MustCmdf(`./bin/golangci-lint run --fix ./...`).RunV()
}

// "Test" Targets

type Test mg.Namespace

// Unit Executes Go unit tests.
func (Test) Unit() error {
	return shx.MustCmdf(`go test -v -count 1 -coverprofile=coverage.out ./...`).RunV()
}

// Coverage Generates file with unit test coverage data and open it in browser
func (t Test) Coverage() error {
	mg.Deps(t.Unit)
	return shx.MustCmdf(`go tool cover -html=coverage.out`).Run()
}
