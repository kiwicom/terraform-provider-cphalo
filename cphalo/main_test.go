package cphalo

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logDebug("EXECUTING TESTMAIN !!!!!!!!")
	setTestID()
	code := m.Run()
	os.Exit(code)
}
