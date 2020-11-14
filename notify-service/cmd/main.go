package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jabardigitalservice/jabarnotify-services/notify-service/pkg/endpoint"
	"github.com/jabardigitalservice/jabarnotify-services/notify-service/pkg/service"
	"github.com/jabardigitalservice/jabarnotify-services/notify-service/pkg/transport"
	"google.golang.org/grpc"
)

type notifyServer struct {
	statsSvcAddr string
	statsSvcConn *grpc.ClientConn
}

func main() {

	var httpAddr = ":8080"

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	addservice := service.NewSiteService(logger)
	addendpoints := endpoint.MakeSiteEndpoints(addservice)
	httpHandlers := transport.MakeHTTPHandler(addendpoints, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", httpAddr)
		server := &http.Server{
			Addr:    httpAddr,
			Handler: httpHandlers,
		}
		errs <- server.ListenAndServe()
	}()

	level.Error(logger).Log("exit", <-errs)
}
