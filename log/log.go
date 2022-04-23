package log

import (
	"bytes"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Formatter struct {
}

func (f *Formatter) Format(e *log.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteString(e.Time.Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString(strings.ToUpper(e.Level.String()))
	b.WriteByte(' ')
	b.WriteString(e.Message)
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func InitLog(level string) {
	logLevels := map[string]log.Level{
		"fatal": log.FatalLevel,
		"error": log.ErrorLevel,
		"warn":  log.WarnLevel,
		"info":  log.InfoLevel,
		"trace": log.TraceLevel,
	}

	if _, ok := logLevels[level]; !ok {
		level = "info"
	}
	log.SetLevel(logLevels[level])
	log.SetFormatter(&Formatter{})
}
