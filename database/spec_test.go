package database

import (
	"github.com/ghthor/gospec"
	"testing"
)

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()

	gospec.MainGoTest(r, t)
}

func TestIntegrationSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeAddFileExecutor)

	gospec.MainGoTest(r, t)
}
