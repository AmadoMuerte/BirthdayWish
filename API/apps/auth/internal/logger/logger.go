package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

type ColorHandler struct {
	handler slog.Handler
	writer  io.Writer
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level
	message := r.Message

	// Выбираем цвет в зависимости от уровня логирования
	var levelColor *color.Color
	switch level {
	case slog.LevelDebug:
		levelColor = color.New(color.FgBlue, color.Bold)
	case slog.LevelInfo:
		levelColor = color.New(color.FgGreen, color.Bold)
	case slog.LevelWarn:
		levelColor = color.New(color.FgYellow, color.Bold)
	case slog.LevelError:
		levelColor = color.New(color.FgRed, color.Bold)
	default:
		levelColor = color.New(color.FgWhite, color.Bold)
	}

	timeStr := r.Time.Format("2006-01-02 15:04:05")

	attrs := ""
	r.Attrs(func(attr slog.Attr) bool {
		if attrs != "" {
			attrs += " "
		}
		attrs += fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any())
		return true
	})

	var logLine string
	if attrs != "" {
		logLine = fmt.Sprintf("%s %s: %s %s\n",
			color.New(color.FgHiBlack).Sprint(timeStr),
			levelColor.Sprint(level.String()),
			message,
			color.New(color.FgHiBlack).Sprint(attrs))
	} else {
		logLine = fmt.Sprintf("%s %s: %s\n",
			color.New(color.FgHiBlack).Sprint(timeStr),
			levelColor.Sprint(level.String()),
			message)
	}

	_, err := h.writer.Write([]byte(logLine))
	return err
}

func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColorHandler{handler: h.handler.WithAttrs(attrs), writer: h.writer}
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {
	return &ColorHandler{handler: h.handler.WithGroup(name), writer: h.writer}
}

func SetupLogger(mode string) *slog.Logger {
	var handler slog.Handler

	if mode == "dev" {
		baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		handler = &ColorHandler{handler: baseHandler, writer: os.Stdout}
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if mode == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
