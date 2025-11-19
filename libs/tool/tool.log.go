package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config holds runtime-configurable values for logging
type LogCfg struct {
	LogDir   string // base directory for logs
	Env      string // e.g. production, development, debug
	MaxSize  int    // lumberjack max size (MB)
	Backups  int    // lumberjack max backups
	MaxAge   int    // lumberjack max age (days)
	Compress bool
	// BodyCaptureLimit controls how many bytes of request/response bodies will be captured.
	// Set to 0 to disable capturing bodies entirely.
	BodyCaptureLimit int
}

// DefaultLogConfig returns sane defaults
func LogConfig() *LogCfg {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "/data/logs"
	}
	env := os.Getenv("GO_ENV")

	return &LogCfg{
		LogDir:           logDir,
		Env:              env,
		MaxSize:          500,   // Maximum log file split size, default 500 MB
		Backups:          50,    // Maximum number of old log files to keep
		MaxAge:           30,    // Maximum number of days to keep old log files
		Compress:         false, // Whether to use gzip to compress and archive log files
		BodyCaptureLimit: 8192,  // 8KB default
	}
}

// Hook implements logrus.Hook
type Hook struct {
	defaultLogger *lumberjack.Logger
	formatter     logrus.Formatter
	loggerByLevel map[logrus.Level]*lumberjack.Logger
}

// Fire writes the formatted entry to the file corresponding to its level.
// It never returns an error to avoid affecting application logic.
func (hook *Hook) Fire(entry *logrus.Entry) error {
	msg, err := hook.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log format error: %v\n", err)
		return nil
	}

	if logger, ok := hook.loggerByLevel[entry.Level]; ok {
		_, werr := logger.Write(msg)
		if werr != nil {
			fmt.Fprintf(os.Stderr, "log write error (level %v): %v\n", entry.Level, werr)
		}
	} else if hook.defaultLogger != nil {
		_, werr := hook.defaultLogger.Write(msg)
		if werr != nil {
			fmt.Fprintf(os.Stderr, "log write error (default): %v\n", werr)
		}
	}

	return nil
}

// Levels returns all log levels because this hook handles level routing internally.
func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// LogInitializer is responsible for building loggers and hooks
type LogInitializer struct {
	cfg *LogCfg
}

func LogFactory(cfg *LogCfg) *LogInitializer {
	if cfg == nil {
		cfg = LogConfig()
	}
	return &LogInitializer{cfg: cfg}
}

func (c *LogInitializer) getLogger(file string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    c.cfg.MaxSize,
		MaxBackups: c.cfg.Backups,
		MaxAge:     c.cfg.MaxAge,
		Compress:   c.cfg.Compress,
		LocalTime:  true,
	}
}

func (c *LogInitializer) createHook(errFile, warFile, infFile, debFile, traFile string) *Hook {
	errlog := c.getLogger(errFile)
	warlog := c.getLogger(warFile)
	inflog := c.getLogger(infFile)
	deblog := c.getLogger(debFile)
	tralog := c.getLogger(traFile)

	hook := Hook{
		defaultLogger: tralog,
		formatter:     &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"},
		loggerByLevel: map[logrus.Level]*lumberjack.Logger{
			logrus.ErrorLevel: errlog,
			logrus.WarnLevel:  warlog,
			logrus.InfoLevel:  inflog,
			logrus.DebugLevel: deblog,
			logrus.TraceLevel: tralog,
		},
	}

	return &hook
}

// HandleLogger configures logrus with hook and level. Call early in main().
func (c *LogInitializer) HandleLog(app string) {
	base := c.cfg.LogDir
	if base == "" {
		base = "/data/logs"
	}

	errFile := filepath.Join(base, app, "error.log")
	warFile := filepath.Join(base, app, "warn.log")
	infFile := filepath.Join(base, app, "info.log")
	debFile := filepath.Join(base, app, "debug.log")
	traFile := filepath.Join(base, app, "trace.log")

	// In debug, prefer local relative logs
	if strings.EqualFold(c.cfg.Env, "debug") {
		pwd, _ := os.Getwd()

		errFile = filepath.Join(pwd, "../logs/error.log")
		warFile = filepath.Join(pwd, "../logs/warn.log")
		infFile = filepath.Join(pwd, "../logs/info.log")
		debFile = filepath.Join(pwd, "../logs/debug.log")
		traFile = filepath.Join(pwd, "../logs/trace.log")
	}

	hook := c.createHook(errFile, warFile, infFile, debFile, traFile)

	// Use discard for the global output — we rely on hook for persistence
	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.AddHook(hook)
}

// logger is a thin wrapper that provides helper log methods and middleware
type AppLogger struct {
	noLineFeed *regexp.Regexp
	cfg        *LogCfg
}

func NewAppLogger(cfg *LogCfg) *AppLogger {
	if cfg == nil {
		cfg = LogConfig()
	}
	return &AppLogger{
		noLineFeed: regexp.MustCompile(`\n|\r|\t`),
		cfg:        cfg,
	}
}

func (l *AppLogger) Error(data any) { logrus.Error(l.encodeJSON(data)) }
func (l *AppLogger) Warn(data any)  { logrus.Warn(l.encodeJSON(data)) }
func (l *AppLogger) Info(data any)  { logrus.Info(l.encodeJSON(data)) }
func (l *AppLogger) Debug(data any) { logrus.Debug(l.encodeJSON(data)) }
func (l *AppLogger) Trace(data any) { logrus.Trace(l.encodeJSON(data)) }

// encodeJSON marshals value to JSON with fallback so that logging never fails silently.
func (l *AppLogger) encodeJSON(v any) string {
	b, err := json.Marshal(v)
	if err == nil {
		return string(b)
	}
	return fmt.Sprintf(`{"marshal_error":"%v","value":"%v"}`, err, v)
}

// // unencodeJSON tries to parse JSON string into map. On failure returns map{"raw": <str>}.
// func (l *AppLogger) unencodeJSON(s string) map[string]any {
// 	var m map[string]any
// 	if err := json.Unmarshal([]byte(s), &m); err != nil {
// 		return map[string]any{"raw": s}
// 	}
// 	return m
// }

// getRequestBodyFromContext attempts to get cached body previously stored in context under key "body".
// If not present, returns nil slice.
func (l *AppLogger) getRequestBodyFromContext(ctx *gin.Context) []byte {
	d, ok := ctx.Get("body")
	if !ok {
		return nil
	}
	if b, ok := d.([]byte); ok {
		return b
	}
	return nil
}

func (l *AppLogger) cleanLineFeed(str string) string {
	return strings.TrimSpace(l.noLineFeed.ReplaceAllString(str, " "))
}

func (l *AppLogger) formatLatency(d time.Duration) string {
	units := []struct {
		size time.Duration
		unit string
	}{
		{time.Minute, "m"},
		{time.Second, "s"},
		{time.Millisecond, "ms"},
		{time.Microsecond, "µs"},
		{time.Nanosecond, "ns"},
	}

	for _, u := range units {
		if d >= u.size {
			if u.size >= time.Second {
				return fmt.Sprintf("%.1f%s", float64(d)/float64(u.size), u.unit)
			}
			return fmt.Sprintf("%d%s", d/u.size, u.unit)
		}
	}
	return "0s"
}

// Trace middleware improvements
// TraceMiddleware returns a middleware that captures request/response metadata safely.
// Behavior:
// - capture request body and response body only up to cfg.BodyCaptureLimit bytes
// - only capture bodies when Content-Type is application/json (configurable easily)
// - avoids buffering for large responses; does not change response behaviors
func (l *AppLogger) TraceMiddleware() gin.HandlerFunc {
	limit := l.cfg.BodyCaptureLimit
	return func(c *gin.Context) {
		start := time.Now()

		var reqBody []byte
		if limit > 0 && c.Request.Body != nil {
			// Use io.ReadAll but limit to prevent OOM
			buf := make([]byte, 0, 512)
			lr := io.LimitReader(c.Request.Body, int64(limit)+1) // read at most limit+1 to detect truncation
			b, _ := io.ReadAll(lr)
			if len(b) > 0 {
				reqBody = b
			}
			// restore body so downstream can read
			c.Request.Body = io.NopCloser(bytes.NewReader(append(buf, reqBody...)))
		}

		// wrap response writer using gin's helper
		bw := &bodyWriter{ResponseWriter: c.Writer, buf: bytes.NewBuffer(nil), limit: limit}
		c.Writer = bw

		// process the request
		c.Next()

		latency := time.Since(start)

		// capture status, uri, ip etc.
		status := c.Writer.Status()
		// only attempt to capture response body if content-type is JSON and limit > 0
		respBodyStr := ""
		if limit > 0 {
			ct := c.Writer.Header().Get("Content-Type")
			if strings.Contains(ct, "application/json") {
				respBytes := bw.buf.Bytes()
				if len(respBytes) > limit {
					respBodyStr = string(respBytes[:limit]) + "...TRUNCATED"
				} else {
					respBodyStr = string(respBytes)
				}
			}
		}

		reqBodyStr := ""
		if len(reqBody) > 0 {
			if len(reqBody) > limit && limit > 0 {
				reqBodyStr = string(reqBody[:limit]) + "...TRUNCATED"
			} else {
				reqBodyStr = string(reqBody)
			}
		}

		l.Trace(map[string]any{
			"uri":       c.Request.RequestURI,
			"method":    c.Request.Method,
			"status":    status,
			"ip":        c.ClientIP(),
			"rip":       c.Request.RemoteAddr,
			"req_id":    c.Request.Header.Get("x-request-id"),
			"req_body":  l.parseJSONOrRaw(l.cleanLineFeed(reqBodyStr)),
			"resp_body": l.parseJSONOrRaw(l.cleanLineFeed(respBodyStr)),
			"dur":       l.formatLatency(latency),
		})
	}
}

// bodyWriter wraps gin.ResponseWriter to capture response bytes up to a limit.
// It preserves all behaviors of the underlying writer.
type bodyWriter struct {
	gin.ResponseWriter
	buf   *bytes.Buffer
	limit int // 0 => disabled
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if w.limit == 0 || w.buf == nil {
		return n, err
	}
	remaining := w.limit - w.buf.Len()
	if remaining <= 0 {
		return n, err
	}
	if len(b) > remaining {
		w.buf.Write(b[:remaining])
	} else {
		w.buf.Write(b)
	}
	return n, err
}

// parseJSONOrRaw tries to unmarshal a string to map; if fails returns raw string under key "raw".
func (l *AppLogger) parseJSONOrRaw(s string) any {
	if s == "" {
		return nil
	}
	var m any
	if json.Unmarshal([]byte(s), &m) == nil {
		return m
	}
	return map[string]any{"raw": s}
}

// LogIllegalEntity logs 4xx illegal entity situations. It reuses getRequestBodyFromContext.
func (l *AppLogger) LogIllegalEntity(c *gin.Context) {
	reqBodyStr := l.cleanLineFeed(string(l.getRequestBodyFromContext(c)))

	l.Warn(map[string]any{
		"uri":      c.Request.RequestURI,
		"method":   c.Request.Method,
		"status":   c.Writer.Status(),
		"ip":       c.ClientIP(),
		"rip":      c.Request.RemoteAddr,
		"req_id":   c.Request.Header.Get("x-request-id"),
		"req_body": l.parseJSONOrRaw(reqBodyStr),
	})
}
