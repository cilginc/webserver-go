package http

import (
	"bufio"
	"errors"
	"strings"
	"unicode"
)

func httpCanonical(s string) string {
	if s == "" {
		return ""
	}
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		parts[i] = string(runes)
	}
	return strings.Join(parts, "-")
}

func readHeaders(r *bufio.Reader) (Header, error) {
	headers := make(Header)

	for {
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}

		if line == "" {
			break
		}

		colon := strings.Index(line, ":")
		if colon == -1 {
			return nil, errors.New("malformed header line: " + line)
		}

		key := httpCanonical(strings.TrimSpace(line[:colon]))
		value := strings.TrimSpace(line[colon+1:])

		headers[key] = append(headers[key], value)
	}

	return headers, nil
}

func readLine(r *bufio.Reader) (string, error) {
	var buf []byte

	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}

		if b == '\r' {
			next, err := r.ReadByte()
			if err != nil {
				return "", err
			}
			if next != '\n' {
				return "", errors.New("expected LF after CR")
			}
			break
		}

		if b == '\n' {
			break
		}

		buf = append(buf, b)
	}

	return string(buf), nil
}
