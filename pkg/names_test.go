package pkg

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SnakeCase", "snake_case"},
		{"Snake_Case", "snake_case"},
		{"snake_case", "snake_case"},
		{"Camel", "camel"},
		{"snake", "snake"},
		{"", ""},
		{"PascalCaseTest", "pascal_case_test"},
		{"HTTPStatusCode", "httpstatus_code"},
		{"commonly.Used", "commonly.used"},
		{"Aaa-Bbb", "aaa_bbb"},
	}

	for _, test := range tests {
		result := ToSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("ToSnakeCase(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
