package templatex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTemplate(t *testing.T) {
	t.Parallel()
	engine := NewTemplateEngine()

	cases := []struct {
		name     string
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "basic_placeholder",
			template: "Hello {{name}}",
			data:     map[string]interface{}{"name": "Alice"},
			expected: "Hello Alice",
		},
		{
			name:     "square_brackets",
			template: "Your code is [[otp]]",
			data:     map[string]interface{}{"otp": "123456"},
			expected: "Your code is 123456",
		},
		{
			name:     "percent_curly",
			template: "Click here: {{link}}",
			data:     map[string]interface{}{"link": "https://example.com"},
			expected: "Click here: https://example.com",
		},
		{
			name:     "fallback_default_value",
			template: "Hi {{name | default:'User'}}",
			data:     map[string]interface{}{},
			expected: "Hi User",
		},
		{
			name:     "fallback_not_used_if_key_exists",
			template: "Hi {{name | default:'User'}}",
			data:     map[string]interface{}{"name": "John"},
			expected: "Hi John",
		},
		{
			name:     "uppercase_filter",
			template: "Hi {{name | uppercase}}",
			data:     map[string]interface{}{"name": "john"},
			expected: "Hi JOHN",
		},
		{
			name:     "nested_key",
			template: "First name: {{user.name.first}}",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": map[string]interface{}{
						"first": "Alice",
					},
				},
			},
			expected: "First name: Alice",
		},
		{
			name:     "mixed_placeholders_with_filters",
			template: "Hi {{name | default:'Friend' | uppercase}}, OTP: [[otp]], Link: {{link | trim}}",
			data: map[string]interface{}{
				"otp":  "999999",
				"link": " https://x.yz ",
			},
			expected: "Hi FRIEND, OTP: 999999, Link: https://x.yz",
		},
		{
			name:     "dollar_braces",
			template: "ID: ${id}",
			data:     map[string]interface{}{"id": "ABC123"},
			expected: "ID: ABC123",
		},
		{
			name:     "double_angle",
			template: "Code: <<code>>",
			data:     map[string]interface{}{"code": "XZY"},
			expected: "Code: XZY",
		},
		{
			name:     "percent_wrapped",
			template: "Token: %%token%%",
			data:     map[string]interface{}{"token": "9999"},
			expected: "Token: 9999",
		},
	}

	for _, tc := range cases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := engine.Render(tc.template, tc.data)
			t.Logf("Template: %s", tc.template)
			t.Logf("Data: %#v", tc.data)
			t.Logf("Expected: %s", tc.expected)
			t.Logf("Got: %s", result)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
