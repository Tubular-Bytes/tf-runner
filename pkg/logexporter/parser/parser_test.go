package parser_test

import (
	"testing"

	"github.com/Tubular-Bytes/tf-runner/pkg/logexporter/parser"
	"github.com/stretchr/testify/require"
)

func TestParseLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		line     []byte
		wantErr  bool
		expected *parser.Entry
	}{
		{
			name:    "valid line",
			line:    []byte(`2025-06-05T08:52:46.760+0100 [DEBUG] using github.com/hashicorp/go-tfe v1.36.0`),
			wantErr: false,
			expected: &parser.Entry{
				Timestamp: "2025-06-05T08:52:46.760+0100",
				Level:     "DEBUG",
				Message:   "using github.com/hashicorp/go-tfe v1.36.0",
			},
		},
		{
			name:    "valid line",
			line:    []byte(`2025-06-05T10:16:37.746+0100 [ERROR] GET http://localhost:3111/state request failed: Get "http://localhost:3111/state": dial tcp [::1]:3111: connect: connection refused`),
			wantErr: false,
			expected: &parser.Entry{
				Timestamp: "2025-06-05T10:16:37.746+0100",
				Level:     "ERROR",
				Message:   `GET http://localhost:3111/state request failed: Get "http://localhost:3111/state": dial tcp [::1]:3111: connect: connection refused`,
			},
		},
		{
			name:    "valid line",
			line:    []byte(`2025-06-05T10:16:38.769+0100 [DEBUG] GET http://localhost:3111/state: retrying in 2s (1 left)`),
			wantErr: false,
			expected: &parser.Entry{
				Timestamp: "2025-06-05T10:16:38.769+0100",
				Level:     "DEBUG",
				Message:   `GET http://localhost:3111/state: retrying in 2s (1 left)`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := parser.ParseTofu(tt.line)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected.Timestamp, output.Timestamp)
			require.Equal(t, tt.expected.Level, output.Level)
			require.Equal(t, tt.expected.Message, output.Message)
		})
	}
}

func TestParseSlog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		line     []byte
		wantErr  bool
		expected *parser.Entry
	}{
		{
			name:    "valid slog line no attributes",
			line:    []byte(`time="2025-06-05T08:52:46.760+0100" level=debug msg="using github.com/hashicorp/go-tfe v1.36.0"`),
			wantErr: false,
			expected: &parser.Entry{
				Timestamp:  "2025-06-05T08:52:46.760+0100",
				Level:      "DEBUG",
				Message:    "using github.com/hashicorp/go-tfe v1.36.0",
				Attributes: nil,
			},
		},
		{
			name:    "valid slog line with attributes",
			line:    []byte(`time="2025-06-05T08:52:46.760+0100" level=debug msg="using github.com/hashicorp/go-tfe v1.36.0" test="foo"`),
			wantErr: false,
			expected: &parser.Entry{
				Timestamp:  "2025-06-05T08:52:46.760+0100",
				Level:      "DEBUG",
				Message:    "using github.com/hashicorp/go-tfe v1.36.0",
				Attributes: map[string]any{"test": "foo"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := parser.ParseSlog(tt.line)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected.Timestamp, output.Timestamp)
			require.Equal(t, tt.expected.Level, output.Level)
			require.Equal(t, tt.expected.Message, output.Message)
			require.Equal(t, tt.expected.Attributes, output.Attributes)
		})
	}
}
