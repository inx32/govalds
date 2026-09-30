package logging

import (
	"io"

	"github.com/rs/zerolog"
)

type zerologSplit struct {
	stdout, stderr io.Writer
	stderrLevel    zerolog.Level
}

func (w *zerologSplit) Write(p []byte) (n int, err error) {
	writer := w.stdout
	if writer == nil {
		writer = w.stderr
	}
	return writer.Write(p)
}

func (w *zerologSplit) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	writer := w.stdout
	if level >= w.stderrLevel {
		writer = w.stderr
	}
	if writer == nil {
		return len(p), nil
	}

	return writer.Write(p)
}

func NewSplit(stdout, stderr io.Writer, stderrLevel zerolog.Level) *zerologSplit {
	if stdout == nil && stderr == nil {
		panic("NewZerologSplit: both writers are nil")
	}
	return &zerologSplit{stdout, stderr, stderrLevel}
}
