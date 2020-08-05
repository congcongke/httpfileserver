package common

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AppErr struct {
	Status  int         `json:"-"`
	Msg     string      `json:"msg"`
	Payload interface{} `json:"payload"`
}

var (
	ENotFound      = &AppErr{http.StatusNotFound, "NotFound", nil}
	EInternalError = &AppErr{http.StatusInternalServerError, "InternalError", nil}
	ETooBusy       = &AppErr{http.StatusConflict, "Too busy", nil}
)

func (e *AppErr) WithPayload(payload interface{}) *AppErr {
	return &AppErr{e.Status, e.Msg, payload}
}

func NewAppErr(code int, msg string, payload interface{}) *AppErr {
	return &AppErr{code, msg, payload}
}

func (e *AppErr) Error() string {
	data, err := json.Marshal(e)
	if err != nil {
		logrus.Errorln("AppErr.Error():", err)
		return ""
	}
	return string(data)
}

func E(c *gin.Context, e error) {
	switch appErr := e.(type) {
	case *AppErr:
		c.JSON(appErr.Status, appErr)
	default:
		c.JSON(EInternalError.Status, EInternalError.WithPayload(e.Error()))
	}
}

func R(c *gin.Context, body interface{}) {
	if c.Request.Method == "POST" {
		c.JSON(http.StatusCreated, body)
		return
	}
	if c.Request.Method == "DELETE" {
		c.JSON(http.StatusNoContent, body)
		return
	}
	c.JSON(http.StatusOK, body)
}
