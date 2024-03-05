package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/nktch1/wearable/internal/wearable_service"
	pbwearable "github.com/nktch1/wearable/pkg/server/wearable"
)

func main() {
	ctx := context.Background()
	address := ":30103"

	if err := Run(ctx, address); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Run(ctx context.Context, address string) error {
	logger := getLogger()
	logger.Info(fmt.Sprintf("Starting listening on %s", address))

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer func() {
		if err = listener.Close(); err != nil {
			logger.Info(fmt.Sprintf("Failed to close server: %s, %s", address, err.Error()))
		}
	}()

	logger.Info(fmt.Sprintf("Listening on %s", address))

	logUnaryInterceptor := grpc.UnaryInterceptor(
		middleware.ChainUnaryServer(
			grpczap.UnaryServerInterceptor(logger),
		),
	)

	server := grpc.NewServer(logUnaryInterceptor)

	service, err := wearable_service.NewService()
	if err != nil {
		return fmt.Errorf("create grpc service: %w", err)
	}

	pbwearable.RegisterWearableServiceServer(server, service)
	reflection.Register(server)

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	return server.Serve(listener)
}

func getLogger() *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		os.Stdout,
		zap.DebugLevel,
	)

	return zap.New(core)
}
