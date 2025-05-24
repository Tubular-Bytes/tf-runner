package logexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"time"

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
	data   []map[string]any
	output io.Writer
	pretty bool
}

func NewLogWriter() *LogWriter {
	return &LogWriter{
		data: make([]map[string]any, 0),
	}
}

func (w *LogWriter) Data() []map[string]any {
	return w.data
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	data := make(map[string]any)
	if err := json.Unmarshal(p, &data); err != nil {
		return 0, err
	}

	w.data = append(w.data, data)

	slices.SortFunc(w.data, func(a, b map[string]any) int {
		if a["@timestamp"] == nil || b["@timestamp"] == nil {
			return 0
		}

		aTime, okA := a["@timestamp"].(string)
		bTime, okB := b["@timestamp"].(string)

		if !okA || !okB {
			return 0
		}

		if aTime < bTime {
			return -1
		} else if aTime > bTime {
			return 1
		}

		return 0
	})

	return len(p), nil
}

func (w *Exporter) Flush(logs []map[string]any, pretty bool) error {
	b := bytes.NewBufferString("")

	encoder := json.NewEncoder(b)
	if pretty {
		encoder.SetIndent("", "  ")
	}

	raw := map[string]any{
		"timestamp":  w.start.Format(time.RFC3339),
		"repository": w.repository,
		"logs":       logs,
	}

	if err := encoder.Encode(raw); err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s.json", w.repository, w.start.Format("20060102150405"))
	slog.Debug("writing logs to object store", "key", key, "size", b.Len())

	_, err := w.client.PutObject(context.Background(), bucket, key, b, int64(b.Len()), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
