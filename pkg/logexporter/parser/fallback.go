package parser

import "time"

func FallbackParser(line []byte) (*Entry, error) {
	if len(line) == 0 {
		return nil, nil // Return nil if the line is empty
	}

	return &Entry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "INFO",
		Message:   string(line),
		Source:    "fallback",
	}, nil
}
