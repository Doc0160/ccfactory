package factory

import (
	"os"
	"time"

	_log "github.com/charmbracelet/log"
)

var log = _log.NewWithOptions(os.Stdout, _log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.DateTime,
	Prefix:          "🏭",
	Level:           _log.DebugLevel,
	CallerFormatter: _log.LongCallerFormatter,
})

var logFactory = _log.NewWithOptions(os.Stdout, _log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.DateTime,
	Prefix:          "📄",
	Level:           _log.DebugLevel,
	CallerFormatter: _log.LongCallerFormatter,
	CallerOffset:    1,
})
