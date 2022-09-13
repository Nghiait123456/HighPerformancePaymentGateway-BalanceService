//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"wire_best_practice/b"
	"wire_best_practice/c"
)

func ProviderB() (b.BInterface, error) {
	wire.Build(
		b.ProviderB,
	)

	return &b.B{}, nil
}

func ProviderC() (c.CInterface, error) {
	wire.Build(
		c.ProviderC,
	)

	return &c.C{}, nil
}
