package secretutil

import "os"

type Provider interface {
	Get(secret string) string
}

type defaultProvider struct{}

func NewProvider() Provider {
	return &defaultProvider{}
}

func (*defaultProvider) Get(secret string) string {
	// TODO: Get from a dedicated secrets' vault
	return os.Getenv(secret)
}
