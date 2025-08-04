package templatex

import (
	"fmt"
	"regexp"
	"strings"
)

type FilterFunc func(input string) string

type TemplateEngine interface {
	Render(tpl string, data map[string]interface{}) (string, error)
	AddFilter(name string, fn FilterFunc)
}

type engine struct {
	filters map[string]FilterFunc
}

type filter struct {
	name string
	args string // chỉ dùng cho default
}

type placeholderDef struct {
	keyIndex       int
	filterRawIndex int
}

var placeholderVariants = []string{
	`\{\{\s*([\w\.-]+)((\s*\|[^}]+)*)\s*\}\}`,  // {{ key | filters }}
	`\[\[\s*([\w\.-]+)((\s*\|[^\]]+)*)\s*\]\]`, // [[ key | filters ]]
	`\$\{\s*([\w\.-]+)((\s*\|[^}]+)*)\s*\}`,    // ${ key | filters }
	`<<\s*([\w\.-]+)((\s*\|[^>]+)*)\s*>>`,      // << key | filters >>
	`%%\s*([\w\.-]+)((\s*\|[^%]+)*)\s*%%`,      // %% key | filters %%
}

var placeholderPattern = regexp.MustCompile(`(?s)` + strings.Join(placeholderVariants, "|"))

var placeholderGroups = []placeholderDef{
	{1, 2},   // {{ }}
	{4, 5},   // [[ ]]
	{7, 8},   // ${}
	{10, 11}, // << >>
	{13, 14}, // %% %%
}

func NewTemplateEngine() TemplateEngine {
	e := &engine{filters: make(map[string]FilterFunc)}
	e.AddFilter("uppercase", strings.ToUpper)
	e.AddFilter("lowercase", strings.ToLower)
	e.AddFilter("trim", strings.TrimSpace)
	capitalize := func(input string) string {
		if input == "" {
			return ""
		}
		return strings.ToUpper(string(input[0])) + input[1:]
	}
	e.AddFilter("capitalize", capitalize)
	e.AddFilter("title", func(s string) string {
		words := strings.Fields(s)
		for i, w := range words {
			if len(w) > 0 {
				words[i] = strings.ToUpper(string(w[0])) + w[1:]
			}
		}
		return strings.Join(words, " ")
	})
	e.AddFilter("titlecase", e.filters["title"])
	e.AddFilter("sentencecase", func(s string) string {
		s = strings.TrimSpace(s)
		if s == "" {
			return ""
		}
		first := strings.ToUpper(s[:1])
		rest := strings.ToLower(s[1:])
		return first + rest
	})
	e.AddFilter("sentence", e.filters["sentencecase"])
	return e
}

func (e *engine) AddFilter(name string, fn FilterFunc) {
	e.filters[name] = fn
}

func (e *engine) Render(tpl string, data map[string]interface{}) (string, error) {
	matches := placeholderPattern.FindAllStringSubmatchIndex(tpl, -1)
	if len(matches) == 0 {
		return tpl, nil
	}

	var b strings.Builder
	last := 0

	for _, m := range matches {
		start, end := m[0], m[1]
		b.WriteString(tpl[last:start])

		var key, rawFilters string

		for _, def := range placeholderGroups {
			ki := def.keyIndex * 2
			fi := def.filterRawIndex * 2

			if ki < len(m) && m[ki] != -1 && m[ki+1] != -1 {
				key = tpl[m[ki]:m[ki+1]]
				if fi < len(m) && m[fi] != -1 && m[fi+1] != -1 && m[fi] < m[fi+1] {
					rawFilters = tpl[m[fi]:m[fi+1]]
				}
				break
			}
		}

		val := lookupValue(data, key)
		filters := parseFilterChain(rawFilters)

		for _, f := range filters {
			if f.name == "default" && val == "" {
				val = f.args
			} else if fn, ok := e.filters[f.name]; ok {
				val = fn(val)
			}
		}

		if val != "" {
			b.WriteString(val)
		} else {
			b.WriteString(tpl[start:end])
		}
		last = end
	}

	b.WriteString(tpl[last:])
	return b.String(), nil
}

func RenderTemplate(tpl string, data map[string]interface{}) string {
	e := NewTemplateEngine()
	result, _ := e.Render(tpl, data)
	return result
}

func parseFilterChain(raw string) []filter {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, "|")
	filters := make([]filter, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.HasPrefix(part, "default:") {
			arg := strings.TrimSpace(part[len("default:"):])
			if len(arg) >= 2 {
				if (arg[0] == '"' && arg[len(arg)-1] == '"') || (arg[0] == '\'' && arg[len(arg)-1] == '\'') {
					arg = arg[1 : len(arg)-1]
				}
			}
			filters = append(filters, filter{name: "default", args: arg})
		} else {
			filters = append(filters, filter{name: part})
		}
	}
	return filters
}

func lookupValue(data map[string]interface{}, path string) string {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current, ok = m[part]
		if !ok {
			return ""
		}
	}

	switch v := current.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case int, int64, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		return ""
	}
}
