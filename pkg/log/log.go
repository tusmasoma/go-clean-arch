package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/slack-go/slack"
)

const (
	NoticeLevel   = 2
	CriticalLevel = 4
)

//nolint:gochecknoglobals // logger is a global variable.
var (
	SeverityDefault  = slog.LevelInfo
	SeverityDebug    = slog.LevelDebug
	SeverityInfo     = slog.LevelInfo
	SeverityNotice   = slog.Level(NoticeLevel)
	SeverityWarning  = slog.LevelWarn
	SeverityError    = slog.LevelError
	SeverityCritical = slog.Level(CriticalLevel)
)

// logger is the global logger.
// it is initialized by init() and should not be modified.
var logger *slog.Logger

// Slack webhook URL and environment.
var (
	environment     = os.Getenv("ENVIRONMENT")
	slackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
)

// init initializes the logger.
//
//nolint:gochecknoinits // init is used for logger initialization.
func init() {
	logFormat := os.Getenv("LOG_FORMAT")
	handler := newHandler(logFormat)
	logger = slog.New(handler)
}

// newHandler returns a slog.Handler based on the given format.
//
//nolint:gocritic // switch is used for future extensibility.
func newHandler(format string) slog.Handler {
	switch format {
	case "json":
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       SeverityDefault,
			ReplaceAttr: attrReplacerForDefault,
		})
	}

	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       SeverityDefault,
		ReplaceAttr: attrReplacerForDefault,
	})
}

// attrReplacerForDefault is default attribute replacer.
func attrReplacerForDefault(_ []string, attr slog.Attr) slog.Attr {
	level, ok := attr.Value.Any().(slog.Level)
	if ok {
		attr.Value = toLogLevel(level)
	}
	return attr
}

// toLogLevel converts a slog.Level to a slog.Value.
//
//nolint:exhaustive // switch is used for future extensibility.
func toLogLevel(level slog.Level) slog.Value {
	var ls string

	switch level {
	case SeverityDebug:
		ls = "DEBUG"
	case SeverityInfo:
		ls = "INFO"
	case SeverityNotice:
		ls = "NOTICE"
	case SeverityWarning:
		ls = "WARNING"
	case SeverityError:
		ls = "ERROR"
	case SeverityCritical:
		ls = "CRITICAL"
	default:
		ls = "DEFAULT"
	}

	return slog.StringValue(ls)
}

// sendSlackNotification sends a message to the configured Slack webhook URL.
func sendSlackNotification(level, msg string, attrs ...any) {
	if environment == "local" {
		return
	}
	if slackWebhookURL == "" {
		return
	}
	text := fmt.Sprintf("[%s] %s", level, msg)
	for _, attr := range attrs {
		text += fmt.Sprintf("\n%v", attr)
	}
	err := slack.PostWebhook(slackWebhookURL, &slack.WebhookMessage{
		Text: text,
	})
	if err != nil {
		Warn("Failed to send Slack notification", slog.Any("error", err))
	}
}

// SetOutput sets the logger output.
func SetOutput(w io.Writer) {
	logger = slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

// Debug logs a debug message.
func Debug(msg string, attrs ...any) {
	DebugContext(context.Background(), msg, attrs...)
}

// DebugContext logs a debug message with a context.
func DebugContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityDebug, msg, attrs...)
}

// Info logs an info message.
func Info(msg string, attrs ...any) {
	InfoContext(context.Background(), msg, attrs...)
}

// InfoContext logs an info message with a context.
func InfoContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityInfo, msg, attrs...)
}

// Notice logs a notice message.
func Notice(msg string, attrs ...any) {
	NoticeContext(context.Background(), msg, attrs...)
}

// NoticeContext logs a notice message with a context.
func NoticeContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityNotice, msg, attrs...)
}

// Warn logs a warning message.
func Warn(msg string, attrs ...any) {
	WarnContext(context.Background(), msg, attrs...)
}

// WarnContext logs a warning message with a context.
func WarnContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityWarning, msg, attrs...)
}

// Error logs an error message.
func Error(msg string, attrs ...any) {
	ErrorContext(context.Background(), msg, attrs...)
	sendSlackNotification("ERROR", msg, attrs...)
}

// ErrorContext logs an error message with a context.
func ErrorContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityError, msg, attrs...)
	sendSlackNotification("ERROR", msg, attrs...)
}

// Critical logs a critical message.
func Critical(msg string, attrs ...any) {
	CriticalContext(context.Background(), msg, attrs...)
	sendSlackNotification("CRITICAL", msg, attrs...)
}

// CriticalContext logs a critical message with a context.
func CriticalContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityCritical, msg, attrs...)
	sendSlackNotification("CRITICAL", msg, attrs...)
}

// Panic logs a critical message and panics.
func Panic(msg string, attrs ...any) {
	PanicContext(context.Background(), msg, attrs...)
	panic(msg)
}

// PanicContext logs a critical message with a context and panics.
func PanicContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityCritical, msg, attrs...)
	sendSlackNotification("CRITICAL", msg, attrs...)
	panic(msg)
}

// Fatal logs a critical message and exits.
func Fatal(msg string, attrs ...any) {
	FatalContext(context.Background(), msg, attrs...)
	os.Exit(1)
}

// FatalContext logs a critical message with a context and exits.
func FatalContext(ctx context.Context, msg string, attrs ...any) {
	logger.Log(ctx, SeverityCritical, msg, attrs...)
	sendSlackNotification("CRITICAL", msg, attrs...)
	os.Exit(1)
}
