package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

const emptyString = ""

func handlePostgresError(c *gin.Context, l *log.Logger, err error, resource string) {
	if err == pg.ErrNoRows {
		c.JSON(http.StatusBadRequest, models.Response{Error: fmt.Sprintf("%s with given id doesn't exist", resource)})
		return
	}
	l.Println(err)

	msg := ""
	pgErr, ok := err.(pg.Error)
	if ok {
		switch pgErr.Field('C') {
		case "23503":
			msg = fmt.Sprintf("%s with given id doesn't exists", resource)
		case "23505":
			msg = fmt.Sprintf("%s with given id already exists", resource)
		}
		if len(msg) > 0 {
			c.JSON(http.StatusBadRequest, models.Response{Error: msg})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
}