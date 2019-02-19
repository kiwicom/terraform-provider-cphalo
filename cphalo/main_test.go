package cphalo

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setTestID()
	code := m.Run()
	os.Exit(code)
}
