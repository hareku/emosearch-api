package main

import (
	"github.com/hareku/emosearch-api/pkg/interfaces/lambda"
	"github.com/hareku/emosearch-api/pkg/registry"
)

func main() {
	registry := registry.NewRegistry()
	handler := lambda.NewLambdaHandler(registry)
	handler.Start()
}
