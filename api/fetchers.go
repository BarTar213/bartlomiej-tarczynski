package api

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/BarTar213/bartlomiej-tarczynski/storage"
	"github.com/BarTar213/bartlomiej-tarczynski/worker"
	"github.com/gin-gonic/gin"
)

const (
	idKey           = "id"
	fetcherResource = "fetcher"

	invalidBodyErr = "invalid body"
	registerJobErr = "couldn't register job associated with fetcher"
)

type FetcherHandlers struct {
	storage     storage.Storage
	worker      *worker.Worker
	fetcherPool *sync.Pool
	logger      *log.Logger
}

func NewFetcherHandlers(s storage.Storage, w *worker.Worker, pool *sync.Pool, l *log.Logger) *FetcherHandlers {
	return &FetcherHandlers{
		storage: s,
		worker:  w,
		fetcherPool: pool,
		logger:  l,
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
	fetcher := h.fetcherPool.Get().(*models.Fetcher)
	defer h.ReturnFetcher(fetcher)

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

	err = h.worker.RegisterJob(fetcher)
	if err != nil {
		h.logger.Printf("Job registration err: %s", err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: registerJobErr})
		return
	}

	err = h.storage.UpdateFetcherJobId(fetcher.Id, fetcher.JobId)
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

	jobId, err := h.storage.DeleteFetcher(id)
	if err != nil {
		handlePostgresError(c, h.logger, err, fetcherResource)
		return
	}
	go h.worker.DeregisterJob(jobId)

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

	jobId, err := h.storage.GetFetcherJob(id)
	if err != nil {
		handlePostgresError(c, h.logger, err, fetcherResource)
		return
	}

	err = h.worker.UpdateJob(fetcher, jobId)
	if err != nil {
		h.logger.Printf("Job update err: %s", err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: registerJobErr})
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
		handlePostgresError(c, h.logger, err, "fetcher history")
		return
	}

	c.JSON(http.StatusOK, history)
}
