package parser

import (
	"strings"
)

type Entry struct {
	Timestamp  string         `json:"timestamp"`
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	Source     string         `json:"source"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

func (e *Entry) HandleLogfmt(key, val []byte) error {
	switch string(key) {
	case "time":
		e.Timestamp = string(val)
	case "level":
		e.Level = strings.ToUpper(string(val))
	case "msg":
		e.Message = string(val)
	default:
		if e.Attributes == nil {
			e.Attributes = make(map[string]any)
		}

		e.Attributes[string(key)] = string(val)
	}

	return nil
}
