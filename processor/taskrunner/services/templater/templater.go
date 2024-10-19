package templater

import (
	"bytes"
	"context"
	"unicode"
	"fmt"
	"regexp"
	"strings"
	"encoding/json"
	hevaluate "ylem_taskrunner/helpers/evaluate"
	"ylem_taskrunner/services/evaluate"
)

var tokenFormat = "$$$__dtmn_tkn#%d"

func ParseTemplate(template string, values interface{}, envVars map[string]interface{}) (string, error) {
	r, err := regexp.Compile("(?s){{(.*?)}}")

	if err != nil {
		return "", err
	}

	template, extractedStringsMap, err := hideQuotedSubstrings(template)

	if err != nil {
		return "", err
	}

	matches := r.FindAllStringSubmatch(template, -1)

	for _, v := range matches {
		ctx := context.WithValue(context.Background(), "ctx", hevaluate.Context{ //nolint:all
			TaskInput: values,
			EnvVars:   envVars,
		})
		result, err := evaluate.AggregateWithContext(ctx, v[1], values)

		if err != nil {
			return "", err
		}

		template = strings.ReplaceAll(template, v[0], fmt.Sprintf("%v", result))
	}

	template = revealQuotedSubstrings(template, extractedStringsMap)

	return template, nil
}

func ParseJsonTemplate(template string, values interface{}, envVars map[string]interface{}) (string, error) {
	whiteSpaceFreeText := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, template)

	if whiteSpaceFreeText == "{{INPUT()}}" {
		return ParseTemplate(template, values, envVars)
	}

	template, extractedStringsMap, err := hideQuotedSubstrings(template)

	if err != nil {
		return "", err
	}

	r, err := regexp.Compile("(?s){{(.*?)}}")

	if err != nil {
		return "", err
	}

	matches := r.FindAllStringSubmatch(template, -1)
	bTemplate := []byte(template)

	for _, v := range matches {
		ctx := context.WithValue(context.Background(), "ctx", hevaluate.Context{ //nolint:all
			TaskInput: values,
			EnvVars:   envVars,
		})
		result, err := evaluate.AggregateWithContext(ctx, v[1], values)

		if err != nil {
			return "", err
		}

		m, err := json.Marshal(result)
		if err != nil {
			return "", err
		}

		bTemplate = bytes.ReplaceAll(bTemplate, []byte(v[0]), m)
	}

	template = revealQuotedSubstrings(string(bTemplate), extractedStringsMap)

	return template, nil
}

func hideQuotedSubstrings(template string) (string, map[string]string, error) {
	r, err := regexp.Compile("\"[^\"]*?(?:{{.*?}})[^\"\\n]*?\"")

	if err != nil {
		return "", nil, err
	}

	extractedStringsMap := map[string]string{}
	matches := r.FindAllStringSubmatch(template, -1)

	for i, v := range matches {
		token := fmt.Sprintf(tokenFormat, i)
		extractedStringsMap[token] = v[0]

		template = strings.ReplaceAll(template, v[0], fmt.Sprintf("%v", token))
	}

	return template, extractedStringsMap, nil
}

func revealQuotedSubstrings(template string, extractedStringsMap map[string]string) string {
	for k, v := range extractedStringsMap {
		template = strings.ReplaceAll(template, k, v)
	}

	return template
}
