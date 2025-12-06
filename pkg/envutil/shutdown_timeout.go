package envutil

import (
	"fmt"
	"time"

	"github.com/apotourlyan/ludus-studii/pkg/envutil/envvar"
)

type shutdownTimeout struct {
	value time.Duration
}

func ShutdownTimeout(provider Provider) Variable[time.Duration] {
	s := provider.Get(envvar.ShutdownTimeout)
	value := parseShutdownTimeout(s)
	return &shutdownTimeout{value}
}

func parseShutdownTimeout(s string) time.Duration {
	if s == "" {
		message := fmt.Sprintf("%q environment var not set", envvar.ShutdownTimeout)
		panic(message)
	}

	timeout, err := time.ParseDuration(s)
	if err != nil {
		message := fmt.Sprintf("%q must be in the format {number}{unit} (30s, 1m, 1m30s)", envvar.ShutdownTimeout)
		panic(message)
	}

	return timeout
}

func (p *shutdownTimeout) Value() time.Duration {
	return p.value
}
