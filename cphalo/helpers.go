package cphalo

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
