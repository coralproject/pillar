package stats

import (
	"os"
	"testing"
)

func setup() {}

func teardown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}
