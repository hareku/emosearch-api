package main

import (
	"github.com/hareku/emosearch-api/pkg/interfaces/lambda/statemachine"
	"github.com/hareku/emosearch-api/pkg/registry"
)

func main() {
	registry := registry.NewRegistry()
	handler := statemachine.NewHandler(registry)
	handler.StartListSearches()
}
