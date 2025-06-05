package parser

import (
	"regexp"

	"github.com/kr/logfmt"
)

func ParseSlog(line []byte) (*Entry, error) {
	if len(line) == 0 {
		return nil, nil // Return nil if the line is empty
	}

	if match, err := regexp.MatchString(`^time=`, string(line)); err != nil || !match {
		return nil, nil // Return nil if the line does not start with "time="
	}

	e := Entry{
		Source: "slog",
	}
	if err := logfmt.Unmarshal(line, &e); err != nil {
		return nil, err
	}

	// e := &Entry{
	// 	Timestamp: raw["time"].(string),                   // Placeholder timestamp
	// 	Level:     strings.ToUpper(raw["level"].(string)), // Placeholder level
	// 	Message:   raw["msg"].(string),                    // Use the line as the message
	// }

	// delete(raw, "time")  // Remove the time field from attributes
	// delete(raw, "level") // Remove the level field from attributes
	// delete(raw, "msg")   // Remove the msg field from attributes
	// e.Attributes = raw   // Assign remaining fields to attributes

	return &e, nil
}
