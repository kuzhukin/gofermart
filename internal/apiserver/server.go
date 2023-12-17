package apiserver

import (
	"context"
	"fmt"
	"gophermart/internal/apiserver/handler"
	"gophermart/internal/apiserver/middleware"
	"gophermart/internal/authservice"
	"gophermart/internal/authservice/authstorage"
	"gophermart/internal/authservice/cryptographer"
	"gophermart/internal/config"
	"gophermart/internal/orderscontroller"
	"gophermart/internal/orderscontroller/ordersstorage"
	"gophermart/internal/sql"
	"gophermart/internal/zlog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type GophermartServer struct {
	srvr    http.Server
	sqlCtrl *sql.Controller

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

	sqlController, err := sql.StartNewController(config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("start new sql controller, err=%w", err)
	}

	if err := registerHandlers(router, sqlController); err != nil {
		return nil, err
	}

	return &GophermartServer{
		srvr: http.Server{
			Addr:    config.RunAddress,
			Handler: router,
		},
		sqlCtrl:           sqlController,
		waitingShutdownCh: make(chan struct{}),
	}, nil
}

func registerMiddlewares(router *chi.Mux) {
	router.Use(middleware.LoggingHTTPHandler)
}

func registerHandlers(router *chi.Mux, sqlCtrl *sql.Controller) error {
	storage := authstorage.New(sqlCtrl)
	cryptographer, err := cryptographer.NewAesCryptographer()
	if err != nil {
		return fmt.Errorf("new aes cryptographer, err=%w", err)
	}

	authService := authservice.NewAuthService(storage, cryptographer)
	ordersController := orderscontroller.NewOrdersController(ordersstorage.New(sqlCtrl))

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
