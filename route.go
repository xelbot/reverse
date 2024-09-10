package reverse

import (
	"fmt"
	"strings"
)

type route struct {
	pattern string
	params  []parameter
}

type parameter struct {
	name        string
	placeholder string
}

func (r *route) url(pairs ...string) (string, error) {
	dict, err := mapFromPairs(pairs...)
	if err != nil {
		return "", err
	}

	url := r.pattern
	for _, param := range r.params {
		if value, found := dict[param.name]; found {
			url = strings.Replace(url, param.placeholder, value, 1)
		}
	}

	return url, nil
}

func createRoute(pattern string) route {
	return route{
		pattern: pattern,
		params:  scanPattern(pattern),
	}
}

func mapFromPairs(pairs ...string) (map[string]string, error) {
	length := len(pairs)
	if length%2 != 0 {
		return nil, fmt.Errorf("reverse: the number of parameters must be even, got %v", pairs)
	}

	dict := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 {
		dict[pairs[i]] = pairs[i+1]
	}

	return dict, nil
}

func scanPattern(pattern string) []parameter {
	var (
		i, cnt int
		buf    []byte
	)

	params := make([]parameter, 0)
	bytes := []byte(pattern)

	length := len(bytes)
	for i = 0; i < length; i++ {
		if bytes[i] == '{' {
			cnt += 1
			if cnt == 1 {
				buf = make([]byte, 0, length)
			}
		}

		if cnt > 0 {
			buf = append(buf, bytes[i])
		}

		if bytes[i] == '}' {
			cnt -= 1
			if cnt == 0 {
				params = append(params, createParameter(string(buf)))
			}
		}
	}

	return params
}

func createParameter(placeholder string) parameter {
	p := parameter{placeholder: placeholder}
	if idx := strings.Index(placeholder, ":"); idx > 0 {
		p.name = placeholder[1:idx]
	} else {
		p.name = placeholder[1 : len(placeholder)-1]
	}

	return p
}
