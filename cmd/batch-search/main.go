package main

import (
	"github.com/hareku/emosearch-api/pkg/interfaces/lambda/batch"
	"github.com/hareku/emosearch-api/pkg/registry"
)

func main() {
	registry := registry.NewRegistry()
	handler := batch.NewLambdaHandler(registry)
	handler.StartSearch()
}
