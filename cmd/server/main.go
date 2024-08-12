package main

import (
	"fmt"
	"metal/internal/pkg/domain/repositories"
	service "metal/internal/server/application/metrics-service"
	"metal/internal/server/presentation/router"
	"net/http"
)

func main() {

	service.SetStorage(repositories.New())
	fmt.Println("Server is up on localhost:8080")
	http.ListenAndServe(":8080", router.Router())
}
