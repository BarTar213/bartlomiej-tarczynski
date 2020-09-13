package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/BarTar213/bartlomiej-tarczynski/storage"
	"github.com/gin-gonic/gin"
)

const (
	idKey           = "id"
	fetcherResource = "fetcher"

	invalidBodyErr = "invalid body"
)

type FetcherHandlers struct {
	storage          storage.Storage
	logger           *log.Logger
}

func NewFetcherHandlers(storage storage.Storage, logger *log.Logger) *FetcherHandlers {
	return &FetcherHandlers{
		storage:          storage,
		logger:           logger,
	}
}

func (h *FetcherHandlers) GetFetchers(c *gin.Context) {
	fetchers, err := h.storage.GetFetchers()
	if err != nil {
		handlePostgresError(c, h.logger, err, fetcherResource)
		return
	}

	c.JSON(http.StatusOK, fetchers)
}

func (h *FetcherHandlers) AddFetcher(c *gin.Context) {
	fetcher := &models.Fetcher{}
	err := c.ShouldBindJSON(fetcher)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidBodyErr})
		return
	}

	if err = fetcher.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	err = h.storage.AddFetcher(fetcher)
	if err != nil {
		handlePostgresError(c, h.logger, err, fetcherResource)
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
		handlePostgresError(c, h.logger, err, fetcherResource)
		return
	}

	c.JSON(http.StatusOK, models.Response{})
}

func (h *FetcherHandlers) UpdateFetcher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(idKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "invalid query param - fetcher id"})
		return
	}

	fetcher := &models.Fetcher{}
	err = c.ShouldBindJSON(fetcher)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidBodyErr})
		return
	}

	if err = fetcher.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	fetcher.Id = id
	err = h.storage.UpdateFetcher(fetcher)
	if err != nil {
		handlePostgresError(c, h.logger, err, fetcherResource)
		return
	}

	c.JSON(http.StatusOK, fetcher)
}

func (h *FetcherHandlers) GetHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(idKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "invalid query param - fetcher id"})
		return
	}

	history, err := h.storage.GetHistory(id)
	if err != nil {
		handlePostgresError(c, h.logger, err, fetcherResource)
		return
	}

	c.JSON(http.StatusOK, history)
}
