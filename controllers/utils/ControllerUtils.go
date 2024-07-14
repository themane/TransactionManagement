package utils

import (
	"net/url"
	"strings"
)

func ParseStrings(urlValues url.Values, paramNames ...string) (map[string]*string, error) {
	parsedParams := map[string]*string{}
	for _, paramName := range paramNames {
		if paramValues, ok := urlValues[paramName]; ok {
			if len(paramValues) == 1 && len(strings.TrimSpace(paramValues[0])) > 0 {
				*parsedParams[paramName] = strings.TrimSpace(paramValues[0])
				continue
			}
		}
		parsedParams[paramName] = nil
	}
	return parsedParams, nil
}
