package cphalo

import (
	"fmt"
	"testing"
)

func Test_AssertStringSlice(t *testing.T) {
	tests := []struct {
		a1, a2 []string
		same   bool
	}{
		{[]string{}, []string{}, true},
		{[]string{""}, []string{"1"}, false},
		{[]string{""}, []string{""}, true},
		{[]string{"1"}, []string{}, false},
		{[]string{"1"}, []string{"1"}, true},
		{[]string{"1", "2"}, []string{"1", "2"}, true},
		{[]string{"1", "2"}, []string{"2", "1"}, false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%.2d", i), func(t *testing.T) {
			if tt.same != assertStringSlice(tt.a1, tt.a2) {
				t.Errorf("expected %v=%v to result in %t", tt.a1, tt.a2, tt.same)
			}
		})
	}
}
