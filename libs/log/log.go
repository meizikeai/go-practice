package log

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func HandleLogger(app string) {
	pwd, _ := os.Getwd()
	mode := os.Getenv("GO_ENV")

	infoPath := filepath.Join("/data/logs/", app, "/info.log")
	errorPath := filepath.Join("/data/logs/", app, "/error.log")

	if mode == "debug" {
		infoPath = pwd + "/logs/info.log"
		errorPath = pwd + "/logs/error.log"
	}

	infoer, _ := rotatelogs.New(
		infoPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(infoPath),
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	errorer, _ := rotatelogs.New(
		errorPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(errorPath),
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writerMap := lfshook.WriterMap{
		logrus.DebugLevel: infoer,
		logrus.InfoLevel:  infoer,
		logrus.WarnLevel:  infoer,
		logrus.ErrorLevel: errorer,
		logrus.FatalLevel: errorer,
		logrus.PanicLevel: errorer,
	}

	logrus.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
}

func AccessLogger(app string) gin.HandlerFunc {
	pwd, _ := os.Getwd()
	mode := os.Getenv("GO_ENV")

	accessPath := filepath.Join("/data/logs/", app, "/access.log")

	if mode == "debug" {
		accessPath = pwd + "/logs/access.log"
	}

	logger := logrus.New()

	accesser, _ := rotatelogs.New(
		accessPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(accessPath),
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	accesserMap := lfshook.WriterMap{
		logrus.InfoLevel:  accesser,
		logrus.FatalLevel: accesser,
		logrus.DebugLevel: accesser,
		logrus.WarnLevel:  accesser,
		logrus.ErrorLevel: accesser,
		logrus.PanicLevel: accesser,
	}

	logger.AddHook(lfshook.NewHook(accesserMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		clientIP := c.ClientIP()
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()

		logger.WithFields(logrus.Fields{
			"client_ip":    clientIP,
			"latency_time": latencyTime,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
			"status_code":  statusCode,
		}).Info()
	}
}
