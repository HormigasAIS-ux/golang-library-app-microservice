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
}
