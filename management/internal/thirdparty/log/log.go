// Package log is a drop-in replacement for the unpublished internal
// chaitin.cn/dev/go/log module. It implements only what this codebase
// actually calls: a small named-logger registry over logrus, an AllLoggers
// sentinel to apply settings (level/formatter/output/hooks) to every
// registered logger at once, and a handful of package-level convenience
// functions (Warn/Infof/Fatal/...) that log through a shared root logger.
package log

import (
	"io"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger is the logger type handed out by GetLogger.
type Logger = logrus.Logger

// TextFormatter is re-exported so callers can configure it directly.
type TextFormatter = logrus.TextFormatter

// AllLoggers is a sentinel value accepted by SetLevel, SetFormatter,
// SetOutput and AddHook meaning "apply to every registered logger".
const AllLoggers = "*"

var (
	mu      sync.Mutex
	root    = logrus.New()
	loggers = map[string]*Logger{}
)

// GetLogger returns the named logger, creating it if necessary. An empty
// name returns the shared root logger, which also backs the package-level
// convenience functions below.
func GetLogger(name string) *Logger {
	mu.Lock()
	defer mu.Unlock()

	if name == "" {
		return root
	}
	if l, ok := loggers[name]; ok {
		return l
	}
	l := logrus.New()
	loggers[name] = l
	return l
}

func resolveTargets(target string) []*Logger {
	if target != AllLoggers {
		return []*Logger{GetLogger(target)}
	}

	mu.Lock()
	defer mu.Unlock()

	all := make([]*Logger, 0, len(loggers)+1)
	all = append(all, root)
	for _, l := range loggers {
		all = append(all, l)
	}
	return all
}

// SetLevel sets the log level on target (a logger name or AllLoggers).
func SetLevel(target string, level logrus.Level) {
	for _, l := range resolveTargets(target) {
		l.SetLevel(level)
	}
}

// SetFormatter sets the formatter on target (a logger name or AllLoggers).
func SetFormatter(target string, formatter logrus.Formatter) {
	for _, l := range resolveTargets(target) {
		l.SetFormatter(formatter)
	}
}

// SetOutput sets the output writer on target (a logger name or AllLoggers).
func SetOutput(target string, out io.Writer) {
	for _, l := range resolveTargets(target) {
		l.SetOutput(out)
	}
}

// AddHook adds a hook on target (a logger name or AllLoggers).
func AddHook(target string, hook logrus.Hook) {
	for _, l := range resolveTargets(target) {
		l.AddHook(hook)
	}
}

// ParseLevel parses a log level name, e.g. "debug", "info", "warn", "error".
func ParseLevel(lvl string) (logrus.Level, error) {
	return logrus.ParseLevel(lvl)
}

// NewLockOutput wraps w so concurrent writes are serialized.
func NewLockOutput(w io.Writer) io.Writer {
	return &lockedWriter{w: w}
}

type lockedWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (lw *lockedWriter) Write(p []byte) (int, error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	return lw.w.Write(p)
}

// NewErrorStackHook returns a hook that, when enabled, attaches a stack
// trace to error-level-and-above log entries.
func NewErrorStackHook(enabled bool) logrus.Hook {
	return &errorStackHook{enabled: enabled}
}

type errorStackHook struct {
	enabled bool
}

func (h *errorStackHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}

func (h *errorStackHook) Fire(entry *logrus.Entry) error {
	if h.enabled {
		entry.Data["stack"] = string(debug.Stack())
	}
	return nil
}

// Package-level convenience functions log through the shared root logger.

func Debug(args ...interface{})                 { root.Debug(args...) }
func Debugf(format string, args ...interface{}) { root.Debugf(format, args...) }
func Info(args ...interface{})                  { root.Info(args...) }
func Infof(format string, args ...interface{})  { root.Infof(format, args...) }
func Infoln(args ...interface{})                { root.Infoln(args...) }
func Warn(args ...interface{})                  { root.Warn(args...) }
func Warnf(format string, args ...interface{})  { root.Warnf(format, args...) }
func Error(args ...interface{})                 { root.Error(args...) }
func Errorf(format string, args ...interface{}) { root.Errorf(format, args...) }
func Fatal(args ...interface{})                 { root.Fatal(args...) }
func Fatalln(args ...interface{})               { root.Fatalln(args...) }
