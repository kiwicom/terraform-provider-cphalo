package cphalo

import (
	"encoding/json"
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

func Test_ExpandStringList(t *testing.T) {
	tests := []struct {
		s []string
	}{
		{[]string{}},
		{[]string{"1"}},
		{[]string{"1", "2", "3"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%.2d", i), func(t *testing.T) {
			s := tt.s
			var i interface{}

			b, err := json.Marshal(s)
			if err != nil {
				t.Fatalf("marhsalling failed: %v", err)
			}

			if err := json.Unmarshal(b, &i); err != nil {
				t.Fatalf("unmarhsalling failed: %v", err)
			}

			got := expandStringList(i)

			if !assertStringSlice(s, got) {
				t.Errorf("expected %v; got %v", s, got)
			}
		})
	}
}
