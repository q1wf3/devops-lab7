package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpctransport "example.com/ticket-service/grpc"
	httptransport "example.com/ticket-service/http"
	"example.com/ticket-service/implementation"
	"example.com/ticket-service/pb"
	"example.com/ticket-service/repository"
	"example.com/ticket-service/transport"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	// 1. Логгер
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "ticket",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	_ = level.Info(logger).Log("msg", "service started")
	defer func() { _ = level.Info(logger).Log("msg", "service ended") }()

	// 2. Сборка слоёв: repository → service → endpoints
	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)

	// 3. HTTP-транспорт
	httpHandler := httptransport.NewService(endpoints, logger)

	// 4. gRPC-транспорт
	grpcServer := grpctransport.NewGRPCServer(endpoints)

	errs := make(chan error)

	// 5. Graceful shutdown: ловим SIGINT/SIGTERM (Ctrl+C)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// 6. HTTP-сервер на :8080
	go func() {
		_ = level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	// 7. gRPC-сервер на :9090
	go func() {
		_ = level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr)
		lis, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		s := grpc.NewServer()
		pb.RegisterTicketServiceServer(s, grpcServer)
		reflection.Register(s) // ← чтобы grpcurl видел сервис без .proto
		errs <- s.Serve(lis)
	}()

	// 8. Блокируем main, пока не придёт сигнал или ошибка
	_ = level.Error(logger).Log("exit", <-errs)
}
