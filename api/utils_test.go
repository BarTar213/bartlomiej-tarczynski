package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BarTar213/bartlomiej-tarczynski/mock"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func Test_handlePostgresError(t *testing.T) {
	type args struct {
		logger   *log.Logger
		err      error
		resource string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "no_rows_error",
			args: args{
				logger:   logger,
				err:      pg.ErrNoRows,
				resource: fetcherResource,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "23503_error",
			args: args{
				logger:   logger,
				err:      &mock.PgError{FieldCode: "23503"},
				resource: fetcherResource,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "23505_error",
			args: args{
				logger:   logger,
				err:      &mock.PgError{FieldCode: "23505"},
				resource: fetcherResource,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not_handled_postgres_error",
			args: args{
				logger:   logger,
				err:      &mock.PgError{FieldCode: "23111"},
				resource: fetcherResource,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)

			handlePostgresError(context, tt.args.logger, tt.args.err, tt.args.resource)

			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}
