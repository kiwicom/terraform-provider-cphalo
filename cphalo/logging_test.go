package cphalo

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	tt := []struct {
		f   func(v ...interface{})
		lvl string
	}{
		{logTrace, "TRACE"},
		{logDebug, "DEBUG"},
		{logInfo, "INFO"},
		{logWarn, "WARN"},
		{logError, "ERROR"},
	}

	for _, test := range tt {
		t.Run(test.lvl, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			msg := "test message"
			test.f(msg)

			got := buf.String()

			expected := fmt.Sprintf("[%s] %s", test.lvl, msg)

			if !strings.Contains(got, expected) {
				t.Errorf("expected `%s` to be substring of `%s`", expected, got)
			}
		})
	}
}

func TestLoggingF(t *testing.T) {
	tt := []struct {
		f   func(format string, v ...interface{})
		lvl string
	}{
		{logTracef, "TRACE"},
		{logDebugf, "DEBUG"},
		{logInfof, "INFO"},
		{logWarnf, "WARN"},
		{logErrorf, "ERROR"},
	}

	for _, test := range tt {
		t.Run(test.lvl, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			format := "test format %d"
			test.f(format, 1)

			got := buf.String()

			expected := fmt.Sprintf("[%s] %s", test.lvl, fmt.Sprintf("test format %d", 1))

			if !strings.Contains(got, expected) {
				t.Errorf("expected `%s` to be substring of `%s`", expected, got)
			}
		})
	}
}
