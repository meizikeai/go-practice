package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Hook struct {
	defaultLogger *lumberjack.Logger
	formatter     logrus.Formatter
	minLevel      logrus.Level
	loggerByLevel map[logrus.Level]*lumberjack.Logger
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	msg, err := hook.formatter.Format(entry)

	if err != nil {
		return err
	}

	if logger, ok := hook.loggerByLevel[entry.Level]; ok {
		_, err = logger.Write([]byte(msg))
	} else {
		_, err = hook.defaultLogger.Write([]byte(msg))
	}

	return err
}

func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.minLevel+1]
}

type CreateLog struct{}

func NewCreateLog() *CreateLog {
	return &CreateLog{}
}

func (c *CreateLog) getLogger(file string) *lumberjack.Logger {
	template := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    100,   // Maximum log file split size, default 100 MB
		MaxBackups: 10,    // Maximum number of old log files to keep
		MaxAge:     15,    // Maximum number of days to keep old log files
		Compress:   false, // Whether to use gzip to compress and archive log files
		LocalTime:  true,  // Whether to use local time, default UTC time
	}

	return template
}

func (c *CreateLog) createHook(errFile, warFile, infFile, debFile, traFile string) *Hook {
	errlog := c.getLogger(errFile)
	warlog := c.getLogger(warFile)
	inflog := c.getLogger(infFile)
	deblog := c.getLogger(debFile)
	tralog := c.getLogger(traFile)

	hook := Hook{
		defaultLogger: tralog,
		minLevel:      logrus.TraceLevel,
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

func (c *CreateLog) HandleLogger(app string) {
	pwd, _ := os.Getwd()
	mode := os.Getenv("GO_ENV")

	errFile := filepath.Join("/data/logs/", app, "/error.log")
	warFile := filepath.Join("/data/logs/", app, "/warn.log")
	infFile := filepath.Join("/data/logs/", app, "/info.log")
	debFile := filepath.Join("/data/logs/", app, "/debug.log")
	traFile := filepath.Join("/data/logs/", app, "/trace.log")

	if mode == "debug" {
		errFile = filepath.Join(pwd, "../logs/error.log")
		warFile = filepath.Join(pwd, "../logs/warn.log")
		infFile = filepath.Join(pwd, "../logs/info.log")
		debFile = filepath.Join(pwd, "../logs/debug.log")
		traFile = filepath.Join(pwd, "../logs/trace.log")
	}

	hook := c.createHook(errFile, warFile, infFile, debFile, traFile)

	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.AddHook(hook)
}

type logger struct {
	noLineFeed *regexp.Regexp
}

func NewLogger() *logger {
	return &logger{
		noLineFeed: regexp.MustCompile(`\n|\r|\t`),
	}
}

func (l *logger) HandleErrorLogging(data any) {
	logrus.Error(l.marshalJson(data))
}

func (l *logger) HandleWarnLogging(data any) {
	logrus.Warn(l.marshalJson(data))
}

func (l *logger) HandleInfoLogging(data any) {
	logrus.Info(l.marshalJson(data))
}

func (l *logger) HandleDebugLogging(data any) {
	logrus.Debug(l.marshalJson(data))
}

func (l *logger) HandleTraceLogging(data any) {
	logrus.Trace(l.marshalJson(data))
}

func (l *logger) marshalJson(date any) string {
	res, err := json.Marshal(date)

	if err != nil {
		return ""
	}

	return string(res)
}

func (l *logger) unmarshalJson(date string) map[string]any {
	res := make(map[string]any, 0)

	_ = json.Unmarshal([]byte(date), &res)

	return res
}

// trace logger
type traceWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (c traceWriter) Write(b []byte) (int, error) {
	c.body.Write(b)
	return c.ResponseWriter.Write(b)
}

func (l *logger) TraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		writer := &traceWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}

		c.Writer = writer

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		body := l.cleanLineFeed(string(l.getMountBody(c)))
		data := writer.body.String()

		l.HandleTraceLogging(map[string]any{
			"uri":     c.Request.RequestURI,
			"method":  c.Request.Method,
			"status":  c.Writer.Status(),
			"client":  c.ClientIP(),
			"remote":  c.Request.RemoteAddr,
			"request": l.getRequestID(c.Request),
			"body":    l.unmarshalJson(body),
			"data":    l.unmarshalJson(data),
			"latency": l.handleTime(latency.Milliseconds()),
		})
	}
}

func (l *logger) LoggingIllegalEntity(c *gin.Context) {
	body := l.cleanLineFeed(string(l.getMountBody(c)))

	l.HandleWarnLogging(map[string]any{
		"error":   "422",
		"uri":     c.Request.RequestURI,
		"method":  c.Request.Method,
		"status":  c.Writer.Status(),
		"client":  c.ClientIP(),
		"remote":  c.Request.RemoteAddr,
		"request": l.getRequestID(c.Request),
		"body":    l.unmarshalJson(body),
	})
}

func (l *logger) getMountBody(ctx *gin.Context) []byte {
	d, _ := ctx.Get("body")
	result, _ := d.([]byte)

	return result
}

func (l *logger) getRequestID(req *http.Request) string {
	rid := req.Header.Get("http_x_request_id")

	if len(rid) != 0 {
		return rid
	}

	return ""
}

func (l *logger) cleanLineFeed(str string) string {
	return l.noLineFeed.ReplaceAllString(str, " ")
}

func (l *logger) handleTime(delta int64) string {
	result := ""

	if delta < 10000 {
		result = fmt.Sprintf("%dms", delta)
	} else {
		result = fmt.Sprintf("%ds", delta/1000)
	}

	return result
}
