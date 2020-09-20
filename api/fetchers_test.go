package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/BarTar213/bartlomiej-tarczynski/config"
	"github.com/BarTar213/bartlomiej-tarczynski/mock"
	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/BarTar213/bartlomiej-tarczynski/storage"
	"github.com/BarTar213/bartlomiej-tarczynski/worker"
	"github.com/gin-gonic/gin"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

const (
	validId     = "5"
	invalidId   = "5a"
	exampleUrl  = "https://httpbin.org/range/15"
	urlKey      = "url"
)

func TestFetcherHandlers_AddFetcher(t *testing.T) {
	type fields struct {
		storage storage.Storage
		logger  *log.Logger
		conf    *config.Config
	}
	tests := []struct {
		name       string
		fields     fields
		body       interface{}
		wantStatus int
	}{
		{
			name: "positive_add_fetcher",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "negative_add_fetcher_content_too_large_error",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			wantStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name: "negative_add_fetcher_invalid_body_error",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: map[string]interface{}{
				urlKey: 65,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_add_fetcher_validation_error",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Interval: 60,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_add_fetcher_update_job_id_error",
			fields: fields{
				storage: &mock.Storage{
					UpdateFetcherJobIdErr: true,
				},
				logger: logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "negative_add_fetcher_storage_error",
			fields: fields{
				storage: &mock.Storage{
					AddFetcherErr: true,
				},
				logger: logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
				WithWorker(),
			)

			jsonBody, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			reqUrl := "/api/fetcher"
			req, _ := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(jsonBody))

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestFetcherHandlers_DeleteFetcher(t *testing.T) {
	type fields struct {
		storage storage.Storage
		logger  *log.Logger
		conf    *config.Config
	}
	tests := []struct {
		name       string
		fields     fields
		fetcherId  string
		wantStatus int
	}{
		{
			name: "positive_delete_fetcher",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			fetcherId:  validId,
			wantStatus: http.StatusOK,
		},
		{
			name: "negative_delete_fetcher_storage_error",
			fields: fields{
				storage: &mock.Storage{
					DeleteFetcherErr: true,
				},
				logger: logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			fetcherId:  validId,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "negative_delete_fetcher_invalid_id",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			fetcherId:  invalidId,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
				WithWorker(),
			)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/api/fetcher/%s", tt.fetcherId)
			req, _ := http.NewRequest(http.MethodDelete, reqUrl, nil)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestFetcherHandlers_GetFetchers(t *testing.T) {
	type fields struct {
		storage storage.Storage
		logger  *log.Logger
		conf    *config.Config
	}
	tests := []struct {
		name       string
		fields     fields
		wantStatus int
	}{
		{
			name: "positive_get_fetchers",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "negative_get_fetchers_storage_error",
			fields: fields{
				storage: &mock.Storage{
					GetFetchersErr: true,
				},
				logger: logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
				WithWorker(),
			)

			w := httptest.NewRecorder()
			reqUrl := "/api/fetcher"
			req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestFetcherHandlers_GetHistory(t *testing.T) {
	type fields struct {
		storage storage.Storage
		logger  *log.Logger
		conf    *config.Config
	}
	tests := []struct {
		name       string
		fields     fields
		fetcherId  string
		wantStatus int
	}{
		{
			name: "positive_get_history",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			fetcherId:  validId,
			wantStatus: http.StatusOK,
		},
		{
			name: "negative_get_history_invalid_id_param_error",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			fetcherId:  invalidId,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_get_history_storage_error",
			fields: fields{
				storage: &mock.Storage{
					GetHistoryErr: true,
				},
				logger: logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			fetcherId:  validId,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
				WithWorker(),
			)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/api/fetcher/%s/history", tt.fetcherId)
			req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestNewFetcherHandlers(t *testing.T) {
	type args struct {
		storage storage.Storage
		logger  *log.Logger
		worker  *worker.Worker
		pool    *sync.Pool
		conf    *config.Config
	}
	tests := []struct {
		name string
		args args
		want *FetcherHandlers
	}{
		{
			name: "positive_new_fetcher_handlers",
			args: args{
				storage: &mock.Storage{},
				logger:  logger,
				conf:    &config.Config{Api: config.Api{MaxContentLength: 1024}},
			},
			want: &FetcherHandlers{
				storage: &mock.Storage{},
				logger:  logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFetcherHandlers(tt.args.storage, tt.args.worker, tt.args.pool, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFetcherHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetcherHandlers_UpdateFetcher(t *testing.T) {
	type fields struct {
		storage storage.Storage
		logger  *log.Logger
		conf    *config.Config
	}
	tests := []struct {
		name       string
		fields     fields
		body       interface{}
		fetcherId  string
		wantStatus int
	}{
		{
			name: "positive_update_fetcher",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			fetcherId:  validId,
			wantStatus: http.StatusOK,
		},
		{
			name: "negative_update_fetcher_invalid_id_error",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			fetcherId:  invalidId,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_update_fetcher_content_too_large_error",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			fetcherId:  validId,
			wantStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name: "negative_update_fetcher_invalid_body_error",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: map[string]interface{}{
				urlKey: 65,
			},
			fetcherId:  validId,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_update_fetcher_validation_error",
			fields: fields{
				storage: &mock.Storage{},
				logger:  logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Interval: 60,
			},
			fetcherId:  validId,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_update_fetcher_get_fetcher_job_error",
			fields: fields{
				storage: &mock.Storage{
					GetFetcherJobErr: true,
				},
				logger: logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			fetcherId:  validId,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "negative_update_fetcher_storage_error",
			fields: fields{
				storage: &mock.Storage{
					UpdateFetcherErr: true,
				},
				logger: logger,
				conf: &config.Config{
					Api: config.Api{MaxContentLength: 1024},
				},
			},
			body: &models.Fetcher{
				Url:      exampleUrl,
				Interval: 60,
			},
			fetcherId:  validId,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
				WithWorker(),
			)

			jsonBody, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/api/fetcher/%s", tt.fetcherId)
			req, _ := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(jsonBody))

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func checkResponseStatusCode(t *testing.T, want int, got int) {
	if want != got {
		t.Errorf("Expected response status code: %d, got: %d", want, got)
	}
}
