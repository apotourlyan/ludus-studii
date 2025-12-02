package timeutil

import "time"

type Provider interface {
	Now() time.Time
}

type DefaultProvider struct{}

func NewProvider() Provider {
	return &DefaultProvider{}
}

func (p *DefaultProvider) Now() time.Time {
	return time.Now()
}
