package logexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Tubular-Bytes/tf-runner/pkg/logexporter/parser"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucket = "terrence-dev"
)

type ExporterConfig struct {
	Repository      string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type Exporter struct {
	client *minio.Client

	repository string
	start      time.Time
}

func New(config ExporterConfig) (*Exporter, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Exporter{
		client:     minioClient,
		repository: config.Repository,
		start:      time.Now(),
	}, nil
}

type LogWriter struct {
	data [][]byte
}

func NewLogWriter() *LogWriter {
	return &LogWriter{
		data: make([][]byte, 0),
	}
}

func (w *LogWriter) Data() [][]byte {
	return w.data
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	if w.data == nil {
		w.data = make([][]byte, 0)
	}

	w.data = append(w.data, p)

	return len(p), nil
}

type SlogWriter struct {
	data []*parser.Entry
}

func (w *Exporter) Flush(logs []byte, pretty bool) error {
	key := fmt.Sprintf("%s/%s.json", w.repository, w.start.Format("20060102150405"))
	slog.Info("flushing logs to store", "endpoint", w.client.EndpointURL(), "bucket", bucket, "object", key)

	buf := bytes.NewBufferString("")

	lines := bytes.Split(logs, []byte("\n"))
	entries := make([]*parser.Entry, 0)

	for _, line := range lines {
		if len(line) == 0 {
			continue // Skip empty lines
		}

		entry, err := parseLine(line)
		if err != nil {
			slog.Warn("failed to parse log line", "line", string(line), "error", err)

			continue // Skip lines that cannot be parsed
		}

		if entry != nil {
			entries = append(entries, entry)
		}
	}

	raw, err := json.Marshal(entries)
	if err != nil {
		slog.Error("failed to marshal log entries", "error", err)

		return err
	}

	buf = bytes.NewBuffer(raw)

	slog.Debug("writing logs to object store", "key", key, "size", buf.Len())

	if _, err := w.client.PutObject(
		context.Background(),
		bucket,
		key,
		buf,
		int64(buf.Len()),
		minio.PutObjectOptions{
			ContentType:     "application/json",
			RetainUntilDate: time.Now().Add(30 * 24 * time.Hour), // Retain for 30 days
			Mode:            minio.Compliance,
			LegalHold:       minio.LegalHoldEnabled,
		}); err != nil {
		return err
	}

	return nil
}

func parseLine(line []byte) (*parser.Entry, error) {
	parsers := []func([]byte) (*parser.Entry, error){
		parser.ParseTofu,
		parser.ParseSlog,
		parser.FallbackParser, // Fallback parser if no other parsers match
	}

	for _, parser := range parsers {
		entry, err := parser(line)
		if err != nil {
			slog.Error("failed to parse log line", "line", string(line), "error", err)

			continue // Try the next parser
		}

		if entry != nil {
			return entry, nil // Return the first successfully parsed entry
		}
	}

	return nil, fmt.Errorf("failed to parse line: %s", string(line))
}
