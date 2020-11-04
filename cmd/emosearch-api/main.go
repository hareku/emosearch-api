package main

import (
	"github.com/hareku/emosearch-api/pkg/interfaces/lambda"
	"github.com/hareku/emosearch-api/pkg/registry"
)

func main() {
	reg := registry.NewRegistry()
	handler := lambda.NewLambdaHandler(reg)
	handler.Start()
}
