package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) registerInstrumentRoutes() {
	s.Router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
