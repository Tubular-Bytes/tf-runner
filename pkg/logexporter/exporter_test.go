package logexporter_test

import (
	"testing"

	"github.com/Tubular-Bytes/tf-runner/pkg/logexporter"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	writer := logexporter.NewLogWriter()
	data := [][]byte{
		[]byte(`{"@timestamp":"2025-05-23T09:32:27.045890+01:00","@level":"info","@message":"Should be last"}`),
		[]byte(`{"@timestamp":"2025-05-21T09:32:27.045890+01:00","@level":"info","@message":"Should be first"}`),
		[]byte(`{"@timestamp":"2025-05-22T09:32:27.045890+01:00","@level":"info","@message":"Test message"}`),
	}

	for _, d := range data {
		n, err := writer.Write(d)
		require.NoError(t, err)
		require.Equal(t, len(d), n)
	}

	logs := writer.Data()

	require.Equal(t, 3, len(logs))
	require.Equal(t, "Should be last", logs[2]["@message"])
	require.Equal(t, "Should be first", logs[0]["@message"])
	require.Equal(t, "Test message", logs[1]["@message"])
}
