package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"

	"github.com/bnb-chain/token-recover-approver/internal/module/approval"
)

type HttpServer struct {
	httpServer      *http.Server
	approvalService *approval.ApprovalService

	logger *zerolog.Logger
}

func NewHttpServer(approvalService *approval.ApprovalService, logger *zerolog.Logger) *HttpServer {
	return &HttpServer{
		approvalService: approvalService,
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

func (server *HttpServer) Shutdown() error {
	return server.httpServer.Close()
}

func (server *HttpServer) setRouter(router *httprouter.Router) {
	server.logger.Info().Msg("http router list")
	server.logger.Info().Msg("GET /ping")
	server.logger.Info().Msg("POST /approve")

	router.GET("/ping", server.Ping)
	router.POST("/approve", server.GetTokenRecoverApproval)
}
