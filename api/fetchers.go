package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/BarTar213/bartlomiej-tarczynski/config"
	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/BarTar213/bartlomiej-tarczynski/storage"
	"github.com/gin-gonic/gin"
)

const (
	idKey = "id"
)

type FetcherHandlers struct {
	storage          storage.Storage
	logger           *log.Logger
	maxContentLength int
}

func NewFetcherHandlers(storage storage.Storage, logger *log.Logger, conf *config.Config) *FetcherHandlers {
	return &FetcherHandlers{
		storage:          storage,
		logger:           logger,
		maxContentLength: conf.Api.MaxContentLength,
	}
}

func (h *FetcherHandlers) GetFetchers(c *gin.Context) {
	fetchers, err := h.storage.GetFetchers()
	if err != nil {
		h.logger.Printf("storage error: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, fetchers)
}

func (h *FetcherHandlers) AddFetcher(c *gin.Context) {
	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "error while reading payload"})
		return
	}

	if len(body) > h.maxContentLength {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, models.Response{Error: "Entity Too Large"})
		return
	}

	fetcher := &models.Fetcher{}
	err = json.Unmarshal(body, fetcher)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "invalid body"})
		return
	}

	if err = fetcher.Validate(); err != nil{
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	err = h.storage.AddFetcher(fetcher)
	if err != nil {
		h.logger.Printf("storage error: %s", err)
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
		h.logger.Printf("storage error: %s", err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
		return
	}

	c.JSON(http.StatusOK, models.Response{})
}

func (h *FetcherHandlers) GetHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(idKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "invalid query param - fetcher id"})
		return
	}

	history, err := h.storage.GetHistory(id)
	if err != nil{
		h.logger.Printf("storage error: %s", err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
		return
	}

	c.JSON(http.StatusOK, history)
}
