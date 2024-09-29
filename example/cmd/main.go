package main

import (
	"context"
	"net/http"

	"github.com/EvgenyOvsov/protoc-gateway-impl/example/api/basic_service"
)

func main() {

	mux := http.NewServeMux()

	basic_service := basic_service.NewBasicService()
	if err := basic_service.Apply(context.Background(), mux); err != nil {
		panic(err)
	}
	print("Listening on :9876")
	print(http.ListenAndServe(":9876", mux))
}
