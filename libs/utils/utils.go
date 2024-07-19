package utils

import (
	"go-practice/libs/tool"

	log "github.com/sirupsen/logrus"
)

func HandleErrorLogging(data any) {
	log.Error(string(tool.MarshalJson(data)))
}

func HandleWarnLogging(data any) {
	log.Warn(string(tool.MarshalJson(data)))
}

func HandleInfoLogging(data any) {
	log.Info(string(tool.MarshalJson(data)))
}

func HandleDebugLogging(data any) {
	log.Debug(string(tool.MarshalJson(data)))
}

func HandleTraceLogging(data any) {
	log.Trace(string(tool.MarshalJson(data)))
}
