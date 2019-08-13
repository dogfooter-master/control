package service

import (
	"dogfooter-control/control/pkg/endpoint"
	"dogfooter-control/control/pkg/grpc"
	"dogfooter-control/control/pkg/grpc/pb"
	"dogfooter-control/control/pkg/http"
	"dogfooter-control/control/pkg/service"
	"flag"
	"fmt"
	stdlog "log"
	"net"
	http1 "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	endpoint1 "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	lightsteptracergo "github.com/lightstep/lightstep-tracer-go"
	"github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
	zipkingoopentracing "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	grpc1 "google.golang.org/grpc"
	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sourcegraph/appdash/opentracing"
)

var tracer opentracinggo.Tracer
var logger log.Logger

var fs = flag.NewFlagSet("control", flag.ExitOnError)
var debugAddr = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
var httpAddr = fs.String("http-addr", ":8081", "HTTP listen address")
var grpcAddr = fs.String("grpc-addr", ":8082", "gRPC listen address")
var thriftAddr = fs.String("thrift-addr", ":8083", "Thrift listen address")
var thriftProtocol = fs.String("thrift-protocol", "binary", "binary, compact, json, simplejson")
var thriftBuffer = fs.Int("thrift-buffer", 0, "0 for unbuffered")
var thriftFramed = fs.Bool("thrift-framed", false, "true to enable framing")
var zipkinURL = fs.String("zipkin-url", "", "Enable Zipkin tracing via a collector URL e.g. http://localhost:9411/api/v1/spans")
var lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
var appdashAddr = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")

func Run(serviceName string) {
	fs.Parse(os.Args[1:])
	*httpAddr = service.GetConfigServerControlHttp()
	*grpcAddr = service.GetConfigServerControlGrpc()

	//logger = log.NewLogfmtLogger(os.Stderr)
	//logger = log.NewLogfmtLogger(log.StdlibWriter{})

	// 윈도우즈 OS에 설치할 경우 사용자로부터 환경변수로 포트를 입력받는다.
	// 환경 변수가 없으면 설정 파일로 ...
	envPort := os.Getenv("DOGFOOTER_CONTROL_PORT")
	if len(envPort) != 0 {
		*httpAddr = ":" + envPort
	}
	if len(serviceName) != 0 {
		logFilePath := os.Getenv("DOGFOOTER_HOME") + "/logs/" + serviceName + "/" + time.Now().Local().Format("20060102150405") + ".log"
		service.CreateDirectory(logFilePath)

		fLog, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			stdlog.Fatalf("Error opening file: %v", err)
		}
		defer fLog.Close()
		logger = log.NewLogfmtLogger(log.NewSyncWriter(fLog))
	} else {
		logger = log.NewLogfmtLogger(os.Stderr)
	}

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	if *zipkinURL != "" {
		logger.Log("tracer", "Zipkin", "URL", *zipkinURL)
		collector, err := zipkingoopentracing.NewHTTPCollector(*zipkinURL)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		defer collector.Close()
		recorder := zipkingoopentracing.NewRecorder(collector, false, "localhost:80", "control")
		tracer, err = zipkingoopentracing.NewTracer(recorder)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
	} else if *lightstepToken != "" {
		logger.Log("tracer", "LightStep")
		tracer = lightsteptracergo.NewTracer(lightsteptracergo.Options{AccessToken: *lightstepToken})
		defer lightsteptracergo.FlushLightStepTracer(tracer)
	} else if *appdashAddr != "" {
		logger.Log("tracer", "Appdash", "addr", *appdashAddr)
		collector := appdash.NewRemoteCollector(*appdashAddr)
		tracer = opentracing.NewTracer(collector)
		defer collector.Close()
	} else {
		logger.Log("tracer", "none")
		tracer = opentracinggo.GlobalTracer()
	}

	svc := service.New(getServiceMiddleware(logger))
	eps := endpoint.New(svc, getEndpointMiddleware(logger))
	g := createService(eps)
	//initMetricsEndpoint(g)
	initCancelInterrupt(g)
	logger.Log("exit", g.Run())

}
func initHttpHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultHttpOptions(logger, tracer)

	httpHandler := http.NewHTTPHandler(endpoints, options)
	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		return http1.Serve(httpListener, httpHandler)
	}, func(error) {
		httpListener.Close()
	})

}
func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}

	return
}
func getEndpointMiddleware(logger log.Logger) (mw map[string][]endpoint1.Middleware) {
	mw = map[string][]endpoint1.Middleware{}
	//addEndpointMiddlewareToAllMethods(mw, endpoint.LoggingMiddleware(logger))
	return
}
func initMetricsEndpoint(g *group.Group) {
	http1.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	debugListener, err := net.Listen("tcp", *debugAddr)
	if err != nil {
		logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
		return http1.Serve(debugListener, http1.DefaultServeMux)
	}, func(error) {
		debugListener.Close()
	})
}
func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

func initGRPCHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultGRPCOptions(logger, tracer)

	grpcServer := grpc.NewGRPCServer(endpoints, options)
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", *grpcAddr)
		baseServer := grpc1.NewServer()
		pb.RegisterControlServer(baseServer, grpcServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})

}
