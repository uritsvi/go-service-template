package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"go-service-template/internal/logger/config"
	"go-service-template/internal/otel"

	"github.com/sirupsen/logrus"
)

var (
	L *logrus.Logger
)

type AlignedTextFormatter struct {
	TimestampFormat string
	MaxMsgLength    int
	AddGoroutineID  bool
}

func (f *AlignedTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := fmt.Sprintf("%-8s", strings.ToUpper(entry.Level.String()))

	msg := entry.Message
	if f.MaxMsgLength > 0 && len(msg) > f.MaxMsgLength {
		msg = msg[:f.MaxMsgLength]
	}

	file := ""
	if entry.Caller != nil {
		file = fmt.Sprintf("%-25s", fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line))
	}

	goroutineID := ""
	if f.AddGoroutineID {
		goroutineID = fmt.Sprintf("%-20s", getGoroutineID())
	}

	log := fmt.Sprintf("%s %s%s%s:: %s\n", timestamp, goroutineID, level, file, msg)
	return []byte(log), nil
}

func getGoroutineID() string {
	buf := make([]byte, 64)
	n := runtime.Stack(buf, false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return fmt.Sprintf("goroutine-%s", idField)
}

func Init(cfg *config.Config) {
	L = logrus.New()
	L.SetReportCaller(true)

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	L.SetLevel(level)

	if cfg.JSONFormat {
		L.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyFile: "filename",
			},
		})
	} else {
		L.SetFormatter(&AlignedTextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			MaxMsgLength:    cfg.MaxMsgLength,
			AddGoroutineID:  cfg.AddGoroutineID,
		})
	}

	if cfg.OtelEnabled {
		err := otel.SetupOtelLogger(cfg.OtelEndpoint, cfg.ServiceName, L)
		if err != nil {
			L.Warnf("Failed to setup OpenTelemetry logger: %v", err)
		}
	}
}
