package deepwildcard

import (
	"context"
	"deepwildcard/internal/deepwildcard/validator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const logTag = "[DW] "

type dhServer struct {
	logger    *log.Logger
	validator *validator.Validator
	config    *Config
	server    struct {
		httpServer        *http.Server
		mux               *http.ServeMux
		stop              chan os.Signal
		stopCtx           context.Context
		stopCtxCancelFunc context.CancelFunc
	}
}

func New(options ...DeepHookOption) (*dhServer, error) {
	var err error
	dw := &dhServer{
		logger: log.Default(),
	}
	for _, o := range options {
		err := o(dw)
		if err != nil {
			return nil, err
		}
	}

	dw.server.mux = http.NewServeMux()
	dw.server.httpServer = &http.Server{
		Addr:    dw.config.ListenAddr,
		Handler: dw.server.mux,
	}
	dw.server.mux.HandleFunc("POST /authenticate/x509", dw.httpX509AuthenticateHandler)

	dw.validator, err = validator.New(validator.WithConfig(&dw.config.ValidatorConfig))
	if err != nil {
		return nil, err
	}

	return dw, nil
}

func (dw *dhServer) ListenAndServe() {

	dw.server.stop = make(chan os.Signal, 1)
	signal.Notify(dw.server.stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		dw.LogF("Starting server on %s\n", dw.server.httpServer.Addr)
		if err := dw.server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			dw.LogF("Server error: %v\n", err)
		}
	}()

	<-dw.server.stop
	dw.Logln("Shutting down server...")

	dw.server.stopCtx, dw.server.stopCtxCancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer dw.server.stopCtxCancelFunc()

	if err := dw.server.httpServer.Shutdown(dw.server.stopCtx); err != nil {
		dw.LogF("Shutdown error: %v\n", err)
	}

	dw.Logln("Server stopped")
}

func (dw *dhServer) LogF(format string, v ...any) {
	format = logTag + format
	dw.logger.Printf(format, v...)
}

func (dw *dhServer) Logln(message string) {
	message = logTag + message
	dw.logger.Println(message)
}
