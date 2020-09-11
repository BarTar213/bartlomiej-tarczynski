package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BarTar213/go-template/models"
	"github.com/BarTar213/go-template/storage"
	"github.com/gin-gonic/gin"
)

const (
	idKey = "id"
)

type FetcherHandlers struct {
	storage storage.Storage
	logger  *log.Logger
}

func NewFetcherHandlers(storage storage.Storage, logger *log.Logger) *FetcherHandlers {
	return &FetcherHandlers{storage: storage, logger: logger}
}

func (h *FetcherHandlers) GetFetchers(c *gin.Context) {
	fetchers, err := h.storage.GetFetchers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, fetchers)
}

func (h *FetcherHandlers) AddFetcher(c *gin.Context) {
	fetcher := &models.Fetcher{}
	err := c.ShouldBindJSON(fetcher)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "invalid body"})
		return
	}

	err = h.storage.AddFetcher(fetcher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
		return
	}
	//todo change return
	c.JSON(http.StatusCreated, fetcher)
}

func (h *FetcherHandlers) DeleteFetcher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(idKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "invalid query param - fetcher id"})
		return
	}

	err = h.storage.DeleteFetcher(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
		return
	}

	c.JSON(http.StatusInternalServerError, models.Response{})
}
