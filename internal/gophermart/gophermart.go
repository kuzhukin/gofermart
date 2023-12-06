package gophermart

import (
	"context"
	"fmt"
	"gophermart/internal/gophermart/authservice"
	"gophermart/internal/gophermart/authservice/cryptographer"
	"gophermart/internal/gophermart/authservice/storage"
	"gophermart/internal/gophermart/config"
	"gophermart/internal/gophermart/handler"
	"gophermart/internal/gophermart/handler/middleware"
	"gophermart/internal/gophermart/orderscontroller"
	"gophermart/internal/gophermart/orderscontroller/ordersstorage"
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

	server, err := newServer(config)
	if err != nil {
		return nil, err
	}
	server.start(config.RunAddress)

	return server, nil
}

func newServer(config *config.Config) (*GophermartServer, error) {
	router := chi.NewRouter()

	registerMiddlewares(router)

	if err := registerHandlers(router); err != nil {
		return nil, err
	}

	return &GophermartServer{
		srvr: http.Server{
			Addr:    config.RunAddress,
			Handler: router,
		},
		waitingShutdownCh: make(chan struct{}),
	}, nil
}

func registerMiddlewares(router *chi.Mux) {
	router.Use(middleware.LoggingHTTPHandler)
}

func registerHandlers(router *chi.Mux) error {
	storage := storage.New()
	cryptographer, err := cryptographer.NewAesCryptographer()
	if err != nil {
		return fmt.Errorf("new aes cryptographer, err=%w", err)
	}

	authService := authservice.NewAuthService(storage, cryptographer)
	ordersController := orderscontroller.NewOrdersController(ordersstorage.New())

	router.Handle(registerEndpoint, handler.NewRegistrationHandler(authService))
	router.Handle(loginEndpoint, handler.NewAutentifiactionHandler(authService))
	router.Handle(ordersEndpoint, handler.NewOrdersHandler(authService, ordersController))
	router.Handle(balanceEndpoint, handler.NewBalanceHandler())
	router.Handle(balanceWithdrawEndpoint, handler.NewBalanceWithdrawHandler())
	router.Handle(allWithdrawalsEndpoint, handler.NewWithdrawalsHandler())

	return nil
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
