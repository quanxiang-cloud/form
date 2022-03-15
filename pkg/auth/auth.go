package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	Auth(*gin.Context) bool

	Filter(*http.Response) error
}
