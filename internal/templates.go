package internal

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type Template struct {
	Source      string
	Destination string
}

func GetTemplatesPaths(allAnswers, values map[string]interface{}, key string) ([]Template, error) {
	val, ok := allAnswers[key]
	if !ok {
		return nil, fmt.Errorf("key %q not found", key)
	}

	interfaces, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("type assertion to []interface{} failed for key %q", key)
	}

	templates := make([]Template, len(interfaces))
	for i, v := range interfaces {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element %d is not a string", i)
		}

		parts := strings.SplitN(str, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid template format in element %d: %q", i, str)
		}

		src, err := renderTemplate(parts[0], values)
		if err != nil {
			return nil, fmt.Errorf("failed to render source in element %d: %w", i, err)
		}

		dest, err := renderTemplate(parts[1], values)
		if err != nil {
			return nil, fmt.Errorf("failed to render destination in element %d: %w", i, err)
		}

		templates[i] = Template{
			Source:      src,
			Destination: dest,
		}
	}

	return templates, nil
}

func renderTemplate(tmplStr string, values map[string]interface{}) (string, error) {
	tmpl, err := template.New("tpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, values); err != nil {
		return "", err
	}

	return buf.String(), nil
}
