package http

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (srv Server) Register() http.Handler {
	router := gin.Default()
	router.Use(cors.Default())

	return router
}
