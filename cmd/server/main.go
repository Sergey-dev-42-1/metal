package main

import (
	service "metal/internal/application/metrics-service"
	"metal/internal/domain/repositories"
	"metal/internal/presentation/router"
	"net/http"
)

func main() {

	service.SetStorage(repositories.New())
	http.ListenAndServe(":8080", router.Router())
}
