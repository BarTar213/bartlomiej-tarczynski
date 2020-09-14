package worker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/BarTar213/bartlomiej-tarczynski/storage"
	"github.com/robfig/cron/v3"
)

type Worker struct {
	c           *cron.Cron
	storage     storage.Storage
	client      *http.Client
	historyPool *sync.Pool
	logger      *log.Logger
}

func New(storage storage.Storage, historyPool *sync.Pool, l *log.Logger) *Worker {
	c := cron.New(cron.WithSeconds())
	c.Start()
	return &Worker{
		c:           c,
		storage:     storage,
		historyPool: historyPool,
		client:      http.DefaultClient,
		logger:      l,
	}
}

func (w *Worker) RegisterJob(fetcher *models.Fetcher) error {
	url := fetcher.Url
	id := fetcher.Id
	entryID, err := w.c.AddFunc(fmt.Sprintf("@every %ds", fetcher.Interval), func() {
		w.processJob(url, id)
	})
	if err != nil {
		return err
	}

	fetcher.JobId = int(entryID)
	return nil
}

func (w *Worker) DeregisterJob(id int) {
	w.c.Remove(cron.EntryID(id))
}

func (w Worker) UpdateJob(fetcher *models.Fetcher, id int) error {
	w.DeregisterJob(id)
	return w.RegisterJob(fetcher)
}

func (w *Worker) processJob(url string, fetcherId int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		w.logger.Printf("NewRequest err: %s", err)
		return
	}

	history := w.historyPool.Get().(*models.History)
	history.Duration = 5
	defer w.ReturnHistoryItem(history)

	t := time.Now()
	response, err := w.client.Do(req)
	duration := time.Since(t).Seconds()
	if err == nil {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err == nil {
			history.Response = pointer(string(body))
		}
	}

	if duration < 5 {
		history.Duration = duration
	}

	history.FetcherId = fetcherId
	history.CreatedAt = t.Unix()

	err = w.storage.AddHistory(history)
	if err != nil {
		w.logger.Printf("AddHistory err: %s", err)
		return
	}
}

func (w *Worker) ReturnHistoryItem(h *models.History) {
	h.Reset()
	w.historyPool.Put(h)
}

func (w *Worker) Stop() {
	ctx := w.c.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func pointer(s string) *string {
	return &s
}
