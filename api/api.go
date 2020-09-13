package api

import (
	"log"

	"github.com/BarTar213/bartlomiej-tarczynski/config"
	"github.com/BarTar213/bartlomiej-tarczynski/middleware"
	"github.com/BarTar213/bartlomiej-tarczynski/storage"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Port    string
	Router  *gin.Engine
	Config  *config.Config
	Storage storage.Storage
	Logger  *log.Logger
}

func WithConfig(conf *config.Config) func(a *Api) {
	return func(a *Api) {
		a.Config = conf
	}
}

func WithLogger(logger *log.Logger) func(a *Api) {
	return func(a *Api) {
		a.Logger = logger
	}
}

func WithStorage(storage storage.Storage) func(a *Api) {
	return func(a *Api) {
		a.Storage = storage
	}
}

func NewApi(options ...func(api *Api)) *Api {
	a := &Api{
		Router: gin.Default(),
	}

	for _, option := range options {
		option(a)
	}

	h := NewFetcherHandlers(a.Storage, a.Logger)

	a.Router.Use(gin.Recovery())
	fetchers := a.Router.Group("/api/fetcher")
	{
		fetchers.GET("", h.GetFetchers)
		fetchers.Use(middleware.CheckContentLength(a.Config.Api.MaxContentLength)).PUT("/:id", h.UpdateFetcher)
		fetchers.Use(middleware.CheckContentLength(a.Config.Api.MaxContentLength)).POST("", h.AddFetcher)
		fetchers.DELETE("/:id", h.DeleteFetcher)

		fetchers.GET("/:id/history", h.GetHistory)
	}

	return a
}

func (a *Api) Run() error {
	return a.Router.Run(a.Config.Api.Port)
}
