package envutil

import "os"

type Provider interface {
	Get(varname string) string
}

type defaultProvider struct{}

func NewProvider() Provider {
	return &defaultProvider{}
}

func (*defaultProvider) Get(varname string) string {
	return os.Getenv(varname)
}
