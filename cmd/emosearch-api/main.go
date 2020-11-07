package main

import (
	"github.com/hareku/emosearch-api/pkg/interfaces/lambda/api"
	"github.com/hareku/emosearch-api/pkg/registry"
)

func main() {
	registry := registry.NewRegistry()
	handler := api.NewLambdaHandler(registry)
	handler.Start()
}
