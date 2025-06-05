package parser

import (
	"regexp"
	"strings"
)

func ParseTofu(line []byte) (*Entry, error) {
	s := strings.TrimSpace(string(line))

	// 2025-06-05T08:52:46.760+0100 [DEBUG] using github.com/hashicorp/go-tfe v1.36.0
	pattern, err := regexp.Compile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}[+-]\d{4}) \[([A-Z]+)\] (.+)$`)
	if err != nil {
		return nil, err // Return error if regex compilation fails
	}

	matches := pattern.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return nil, nil // No matches found, return nil
	}

	return &Entry{
		Timestamp: matches[0][1],
		Level:     matches[0][2],
		Message:   matches[0][3],
		Source:    "tofu",
	}, nil
}
