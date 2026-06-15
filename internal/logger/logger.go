package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

const (
	Reset        = "\033[0m"
	ColorInfo    = "\033[34m"
	ColorWarn    = "\033[33m"
	ColorError   = "\033[31m"
	ColorDebug   = "\033[36m"
	ColorDefault = "\033[90m"
)

type ColorHandler struct {
	out   io.Writer
	level slog.Level
}

func NewColorHandler(out io.Writer, level slog.Level) *ColorHandler {
	return &ColorHandler{out: out, level: level}
}

func (h *ColorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *ColorHandler) Handle(_ context.Context, r slog.Record) error {
	var colorCode, msgType string

	switch r.Level {
	case slog.LevelInfo:
		colorCode = ColorInfo
		msgType = "[INFO]"
	case slog.LevelError:
		colorCode = ColorError
		msgType = "[ERROR]"
	case slog.LevelWarn:
		colorCode = ColorWarn
		msgType = "[WARN]"
	case slog.LevelDebug:
		colorCode = ColorDebug
		msgType = "[DEBUG]"
	default:
		colorCode = ColorDefault
		msgType = ""
	}

	colorizedLabel := colorCode + msgType + Reset + ColorDefault + " " + r.Time.Format("15:04:05") + " " + Reset + r.Message
	fmt.Fprintf(h.out, "%s", colorizedLabel)

	fmt.Fprintln(h.out)
	return nil
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {
	return h
}

func CustomLogger(level slog.Level) *slog.Logger {
	customHandler := NewColorHandler(os.Stdout, level)
	return slog.New(customHandler)
}
