package gophermart

import (
	"context"
	"gophermart/internal/gophermart/config"
	"gophermart/internal/gophermart/handler"
	"gophermart/internal/gophermart/handler/middleware"
	"gophermart/internal/gophermart/zlog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type GophermartServer struct {
	srvr http.Server

	waitingShutdownCh chan struct{}
}

func StartNew() (*GophermartServer, error) {
	config, err := config.Make()
	if err != nil {
		return nil, err
	}

	server := newServer(config)
	server.start(config.RunAddress)

	return server, nil
}

func newServer(config *config.Config) *GophermartServer {
	router := chi.NewRouter()

	registerMiddlewares(router)
	registerHAndlers(router)

	return &GophermartServer{
		srvr: http.Server{
			Addr:    config.RunAddress,
			Handler: router,
		},
		waitingShutdownCh: make(chan struct{}),
	}
}

func registerMiddlewares(router *chi.Mux) {
	router.Use(middleware.LoggingHTTPHandler)
}

func registerHAndlers(router *chi.Mux) {
	router.Handle(registerEndpoint, handler.NewRegisterHandler())
	router.Handle(loginEndpoint, handler.NewLoginHandler())
	router.Handle(ordersEndpoint, handler.NewOrdersHandler())
	router.Handle(balanceEndpoint, handler.NewBalanceHandler())
	router.Handle(balanceWithdrawEndpoint, handler.NewBalanceWithdrawHandler())
	router.Handle(allWithdrawalsEndpoint, handler.NewWithdrawalsHandler())
}

func (s *GophermartServer) start(hostport string) {
	go func() {
		defer close(s.waitingShutdownCh)

		if err := s.srvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Logger.Errorf("Http listen and serve address=%s, err=%s", s.srvr.Addr, err)
		}
	}()
}

func (s *GophermartServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return s.srvr.Shutdown(ctx)
}

func (s *GophermartServer) Wait() <-chan struct{} {
	return s.waitingShutdownCh
}
