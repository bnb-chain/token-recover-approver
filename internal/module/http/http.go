package http

import (
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/bnb-chain/token-recover-approver/internal/module/approval"
)

type HttpServer struct {
	httpServer      *http.Server
	approvalService *approval.ApprovalService

	registry      *prometheus.Registry
	metricsServer *http.Server

	logger *zerolog.Logger
}

func NewHttpServer(approvalService *approval.ApprovalService, registry *prometheus.Registry, logger *zerolog.Logger) *HttpServer {
	return &HttpServer{
		approvalService: approvalService,
		registry:        registry,
		logger:          logger,
	}
}

func (server *HttpServer) Run(addr string) error {
	router := httprouter.New()
	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}
	server.setRouter(router)
	return server.httpServer.ListenAndServe()
}

func (server *HttpServer) RunMetrics(addr string, metricsPath string, enablePProf bool) error {
	router := httprouter.New()
	server.metricsServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}
	server.setMetrics(router, metricsPath, enablePProf)
	return server.metricsServer.ListenAndServe()
}

func (server *HttpServer) Shutdown() error {
	err := server.httpServer.Close()
	if err != nil {
		return err
	}

	if server.metricsServer != nil {
		err := server.metricsServer.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (server *HttpServer) setRouter(router *httprouter.Router) {
	server.logger.Info().Msg("http router list")
	server.logger.Info().Msg("GET /ping")
	server.logger.Info().Msg("POST /approve")

	router.GET("/ping", server.Ping)
	router.POST("/approve", server.GetTokenRecoverApproval)
}

func (server *HttpServer) setMetrics(router *httprouter.Router, path string, enablePProf bool) {
	server.logger.Info().Msg("metrics router list")
	server.logger.Info().Msgf("GET %s", path)

	router.GET(path, wrapHttpHandler(promhttp.HandlerFor(server.registry, promhttp.HandlerOpts{})))

	if enablePProf {
		server.logger.Info().Msg("pprof router list")
		server.logger.Info().Msg("GET /debug/pprof/")
		server.logger.Info().Msg("GET /debug/pprof/cmdline")
		server.logger.Info().Msg("GET /debug/pprof/profile")
		server.logger.Info().Msg("GET /debug/pprof/symbol")
		server.logger.Info().Msg("GET /debug/pprof/trace")
		server.logger.Info().Msg("GET /debug/pprof/goroutine")
		server.logger.Info().Msg("GET /debug/pprof/heap")
		server.logger.Info().Msg("GET /debug/pprof/allocs")
		server.logger.Info().Msg("GET /debug/pprof/threadcreate")
		server.logger.Info().Msg("GET /debug/pprof/block")
		server.logger.Info().Msg("GET /debug/pprof/mutex")
		server.logger.Info().Msg("GET /debug/fgprof")

		router.GET("/debug/pprof/", wrapHttpHandleFunc(pprof.Index))
		router.GET("/debug/pprof/cmdline", wrapHttpHandleFunc(pprof.Cmdline))
		router.GET("/debug/pprof/profile", wrapHttpHandleFunc(pprof.Profile))
		router.GET("/debug/pprof/symbol", wrapHttpHandleFunc(pprof.Symbol))
		router.GET("/debug/pprof/trace", wrapHttpHandleFunc(pprof.Trace))
		router.GET("/debug/pprof/goroutine", wrapHttpHandler(pprof.Handler("goroutine")))
		router.GET("/debug/pprof/heap", wrapHttpHandler(pprof.Handler("heap")))
		router.GET("/debug/pprof/allocs", wrapHttpHandler(pprof.Handler("allocs")))
		router.GET("/debug/pprof/threadcreate", wrapHttpHandler(pprof.Handler("threadcreate")))
		router.GET("/debug/pprof/block", wrapHttpHandler(pprof.Handler("block")))
		router.GET("/debug/pprof/mutex", wrapHttpHandler(pprof.Handler("mutex")))
		router.GET("/debug/fgprof", wrapHttpHandler(fgprof.Handler()))
	}
}

// wrapHttpHandler transforms ususal handlers (http.Handler) of the standard library
// package net/http into valid ones of httprouter by adding the params
// (httprouter.Params) parameter.
//
// Use this function to add a context to middlewares or other things that should be
// shared between both handler types.
func wrapHttpHandler(h http.Handler) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		context.Set(req, "params", ps)
		h.ServeHTTP(rw, req)
	}
}

// wrapHttpHandler transforms ususal handlers (http.HandlerFunc) of the standard library
// package net/http into valid ones of httprouter by adding the params
// (httprouter.Params) parameter.
//
// Use this function to add a context to middlewares or other things that should be
// shared between both handler types.
func wrapHttpHandleFunc(h http.HandlerFunc) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		context.Set(req, "params", ps)
		h(rw, req)
	}
}
