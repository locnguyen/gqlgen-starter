package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"gqlgen-starter/cmd/build"
	"gqlgen-starter/config"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

var onceLogger sync.Once
var logger zerolog.Logger

// GetLogger reference: https://betterstack.com/community/guides/logging/zerolog/
func GetLogger() zerolog.Logger {
	onceLogger.Do(func() {
		logLevel, err := zerolog.ParseLevel(config.Application.LogLevel)
		if err != nil {
			logLevel = zerolog.InfoLevel
		}

		zerolog.LevelFieldName = "severity"
		// Usually we'll always want logs as JSON unless we're working on our local machine
		if config.Application.StructuredLogging {
			logger = zerolog.New(osWriter{Out: os.Stdout}).
				Level(logLevel).
				With().
				Timestamp().
				Str("BuildCommit", build.BuildCommit).
				Str("BuildVersion", build.BuildVersion).
				Logger()
		} else {
			logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
				Level(logLevel).
				With().
				Timestamp().
				Caller().
				Logger()
		}
	})
	return logger
}

type osWriter struct {
	Out io.Writer
}

// Write is from the io.Writer interface.  Logic taken partially from zerolog.ConsoleWriter
// Parses top level json fields. If value is numeric, converts to string and writes out to buffer.
// e.g. `{"foo":1234567890}` --> `{"foo":"1234567890"}`
func (t osWriter) Write(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	defer func() {
		buf.Reset()
	}()

	// parse json, if field type is int64, convert to string, write to buffer
	var evt map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&evt)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	var fields = make([]string, 0, len(evt))
	for field := range evt {
		fields = append(fields, field)
	}

	buf.WriteByte('{')

	for i, field := range fields {
		buf.WriteString(strconv.Quote(field) + ":")
		switch fValue := evt[field].(type) {
		case json.Number:
			buf.WriteString(strconv.Quote(fValue.String()))
		default:
			b, err := json.Marshal(fValue)
			if err != nil {
				buf.WriteString(fmt.Sprintf("%v", err))
			} else {
				buf.Write(b)
			}
		}

		if i < len(fields)-1 { // Skip comma for last field
			buf.WriteByte(',')
		}
	}

	buf.WriteByte('}')
	buf.WriteByte('\n')

	numBytes, err := buf.WriteTo(t.Out)
	return int(numBytes), err
}
