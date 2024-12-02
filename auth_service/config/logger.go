package config

import (
	"os"

	"github.com/op/go-logging"
)

func ConfigureLogger() {
	format := logging.MustStringFormatter(
		`%{color}[%{time:2006-01-02 15:04:05]} ▶ %{level}%{color:reset} %{message} ...[%{shortfile}@%{shortfunc}()]`,
	)

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	logging.SetBackend(backend2Formatter)
	switch Envs.LOG_LEVEL {
	case "debug":
		logging.SetLevel(logging.DEBUG, "main")
	case "info":
		logging.SetLevel(logging.INFO, "main")
	case "warning":
		logging.SetLevel(logging.WARNING, "main")
	case "error":
		logging.SetLevel(logging.ERROR, "main")
	}
}
