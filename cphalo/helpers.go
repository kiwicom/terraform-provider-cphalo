package cphalo

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"text/template"
	"time"
)

var testID string

func assertStringSlice(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func expandStringList(input interface{}) []string {
	interSlice := input.([]interface{})
	vs := make([]string, 0, len(interSlice))

	for _, v := range interSlice {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}

	return vs
}

func isCI() bool {
	return len(os.Getenv("CI")) > 0
}

func readTestTemplateData(filePath, uniqueID string) (string, error) {
	path := fmt.Sprintf("testdata/%s", filePath)

	var b bytes.Buffer
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("could not parse template file %s: %v", path, err)
	}

	err = tmpl.Execute(&b, struct {
		Prefix string
	}{
		Prefix: uniqueID,
	})
	if err != nil {
		return "", fmt.Errorf("cannot read execute template from file %s: %v", path, err)
	}

	return b.String(), nil
}

func setTestID() {
	testID = os.Getenv("CPHALO_TEST_ID")
	if testID == "" {
		h := sha1.New()
		_, _ = h.Write([]byte(time.Now().String() + string(rand.Int())))
		testID = fmt.Sprintf("%x", h.Sum(nil))[:6]
	}

	testID += "_"
}
