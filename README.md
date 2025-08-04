# templatex
Template Engine (templatex)

A lightweight and extensible Go template engine that supports multiple placeholder formats, nested key lookup, and flexible filters.

## Features

* Support for various placeholder syntaxes:

  * `{{ key }}`
  * `[[key]]`
  * `${key}`
  * `<<key>>`
  * `%%key%%`
* Nested key resolution: `{{ user.name.first }}`
* Filter support:

  * `default:'fallback'`
  * `uppercase`, `lowercase`, `trim`
  * `titlecase`, `sentencecase`, `uppercaseFirst`
* Custom filter registration

---

## Installation

```bash
go get github.com/vuongthai91/templatex
```

---

## Usage

```go
package main

import (
	"fmt"
	"github.com/your-org/templatex"
)

func main() {
	engine := templatex.NewTemplateEngine()

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": map[string]interface{}{
				"first": "alice",
			},
		},
		"otp":  "123456",
		"link": " https://example.com ",
	}

	tpl := "Hi {{user.name.first | uppercaseFirst}}, OTP: [[otp]], Link: ${link|trim}"
	result, _ := engine.Render(tpl, data)
	fmt.Println(result)
	// Output: Hi Alice, OTP: 123456, Link: https://example.com
}
```

---

## Available Filters

| Filter           | Description                                   |
| ---------------- | --------------------------------------------- |
| `default`        | Fallback if the key is empty                  |
| `uppercase`      | Converts to UPPERCASE                         |
| `lowercase`      | Converts to lowercase                         |
| `trim`           | Trims spaces from both ends                   |
| `uppercaseFirst` | Capitalizes the first letter only             |
| `titlecase`      | Capitalizes the first letter of every word    |
| `sentencecase`   | Capitalizes only the first letter of sentence |

You can chain filters:

```gotemplate
{{ name | default:'Friend' | uppercase }}
```

---

## Custom Filters

You can register your own filters:

```go
engine.AddFilter("reverse", func(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
})
```

---

## Tests

Run unit tests:

```bash
go test ./...
```

