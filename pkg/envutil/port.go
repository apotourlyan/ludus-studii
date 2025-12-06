package envutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apotourlyan/ludus-studii/pkg/envutil/envvar"
)

type port struct {
	value int
}

func Port(provider Provider) Variable[int] {
	portstr := provider.Get(envvar.Port)
	value := parsePort(portstr)
	return &port{value}
}

func parsePort(s string) int {
	if s == "" {
		message := fmt.Sprintf("%q environment var not set", envvar.Port)
		panic(message)
	}

	s, _ = strings.CutPrefix(s, ":")

	port, err := strconv.Atoi(s)
	if err != nil || port < 1024 || port > 65535 {
		message := fmt.Sprintf("%q must be an integer in the range (1024-65535)", envvar.Port)
		panic(message)
	}

	return port
}

func (p *port) Value() int {
	return p.value
}
