package secretutil

import (
	"github.com/apotourlyan/ludus-studii/pkg/secretutil/secretvar"
)

type connectionString struct {
	value string
}

func ConnectionString(provider Provider) Secret[string] {
	s := provider.Get(secretvar.DbConnection)
	validateConnectionString(s)
	return &connectionString{s}
}

func validateConnectionString(s string) {
	// TODO: validate
}

func (p *connectionString) Value() string {
	return p.value
}
