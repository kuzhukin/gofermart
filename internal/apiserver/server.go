package apiserver

import (
	"context"
	"fmt"
	"gophermart/internal/apiserver/handler"
	"gophermart/internal/apiserver/middleware"
	"gophermart/internal/authservice"
	"gophermart/internal/authservice/cryptographer"
	"gophermart/internal/balancecontroller"
	"gophermart/internal/config"
	"gophermart/internal/orderscontroller"
	"gophermart/internal/orderscontroller/accrual"
	"gophermart/internal/storage/accrualstorage"
	"gophermart/internal/storage/ordersstorage"
	"gophermart/internal/storage/sql"
	"gophermart/internal/storage/userstorage"
	"gophermart/internal/zlog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type GophermartServer struct {
	srvr http.Server

	sqlCtrl     *sql.Controller
	accrualCtrl *accrual.AccrualController

	authService *authservice.AuthService
	ordersCtrl  *orderscontroller.OrdersController
	balanceCtrl *balancecontroller.Controller

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

	sqlController, err := sql.StartNewController(config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("start new sql controller, err=%w", err)
	}

	userStorage := userstorage.New(sqlController)
	orderStorage := ordersstorage.New(sqlController)
	accrualStorage := accrualstorage.New(sqlController)

	cryptographer, err := cryptographer.NewAesCryptographer()
	if err != nil {
		return nil, fmt.Errorf("new aes cryptographer, err=%w", err)
	}

	authService := authservice.NewAuthService(userStorage, cryptographer)
	ordersCtrl := orderscontroller.NewOrdersController(orderStorage)
	balancecCtrl := balancecontroller.New(userStorage)

	accrualCtrl := accrual.StartNewController(orderStorage, accrualStorage, config.AccrualAddress)

	server := &GophermartServer{
		sqlCtrl:           sqlController,
		accrualCtrl:       accrualCtrl,
		authService:       authService,
		ordersCtrl:        ordersCtrl,
		balanceCtrl:       balancecCtrl,
		waitingShutdownCh: make(chan struct{}),
	}

	server.initHTTPServer(config.RunAddress)

	return server, nil

}

func (s *GophermartServer) initHTTPServer(addr string) {
	router := chi.NewRouter()

	router.Use(middleware.LoggingHTTPHandler)

	router.Handle(registerEndpoint, handler.NewRegistrationHandler(s.authService))
	router.Handle(loginEndpoint, handler.NewAutentifiactionHandler(s.authService))
	router.Handle(ordersEndpoint, handler.NewOrdersHandler(s.authService, s.ordersCtrl))
	router.Handle(balanceEndpoint, handler.NewBalanceHandler(s.authService, s.balanceCtrl))
	router.Handle(balanceWithdrawEndpoint, handler.NewBalanceWithdrawHandler())
	router.Handle(allWithdrawalsEndpoint, handler.NewWithdrawalsHandler())

	s.srvr = http.Server{Addr: addr, Handler: router}
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
