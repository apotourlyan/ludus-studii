package main

import (
	"github.com/apotourlyan/ludus-studii/pkg/envutil"
	"github.com/apotourlyan/ludus-studii/pkg/secretutil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/app"
)

func main() {
	ep := envutil.NewProvider()
	sp := secretutil.NewProvider()
	app.New(ep, sp).Run()
}
