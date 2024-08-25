package main

import (
	"metal/internal/pkg/domain/repositories"
	service "metal/internal/server/application/metrics-service"
	"metal/internal/server/presentation/controller"
	"metal/internal/server/presentation/router"
)

func main() {
	parseFlags()

	r := router.Router()
	r = controller.AddMetricRoutes(r)
	//Не будет работать если запускать сервер не из корневой папки
	r.LoadHTMLGlob("internal/server/presentation/templates/*.html")
	service.SetStorage(repositories.New())

	r.Run(startAddress.String())
}
