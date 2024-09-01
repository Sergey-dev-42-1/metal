package main

import (
	"metal/internal/pkg/domain/repositories"
	"metal/internal/pkg/gzip"
	"metal/internal/pkg/logger"
	service "metal/internal/server/application/metrics-service"
	"metal/internal/server/presentation/controller"
	"metal/internal/server/presentation/router"
)

func main() {
	parseFlags()
	r := router.Router()
	logger.New("debug")
	r.Use(logger.Logger())
	r.Use(gzip.GzipHandler())

	mc := controller.New(r)
	r = mc.AddRoutes()
	//Не будет работать если запускать сервер не из корневой папки
	r.LoadHTMLGlob("internal/server/presentation/templates/*.html")
	service.SetStorage(repositories.New())
	r.Run(startAddress.String())

}
