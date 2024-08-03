package utils

import (
	"go-practice/libs/tool"

	log "github.com/sirupsen/logrus"
)

var units = tool.NewUnits()
var rules = tool.NewRegexp()

func HandleErrorLogging(data any) {
	log.Error(string(units.MarshalJson(data)))
}

func HandleWarnLogging(data any) {
	log.Warn(string(units.MarshalJson(data)))
}

func HandleInfoLogging(data any) {
	log.Info(string(units.MarshalJson(data)))
}

func HandleDebugLogging(data any) {
	log.Debug(string(units.MarshalJson(data)))
}

func HandleTraceLogging(data any) {
	log.Trace(string(units.MarshalJson(data)))
}
