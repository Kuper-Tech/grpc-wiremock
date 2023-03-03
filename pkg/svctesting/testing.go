package svctesting

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	healther "github.com/SberMarket-Tech/grpc-wiremock/pkg/svctesting/client"
)

var (
	healthRequestTimeout  = time.Second * 50
	healthRequestInterval = time.Second
)

type tester struct {
	interval time.Duration
	client   healther.HealthClient
}

func createClient(port string) (healther.HealthClient, error) {
	addr := fmt.Sprintf(":%s", port)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create health connect: %w", err)
	}

	return healther.NewHealthClient(conn), nil
}

func Check(ctx context.Context, port string) error {
	ctx, cancel := context.WithTimeout(ctx, healthRequestTimeout)
	defer cancel()

	healthClient, err := createClient(port)
	if err != nil {
		return fmt.Errorf("create health client: %w", err)
	}

	t := tester{
		interval: healthRequestInterval,
		client:   healthClient,
	}

	if err = t.isHealthy(ctx); err != nil {
		return fmt.Errorf("is service healthy: %w", err)
	}

	return nil
}

func (t *tester) isHealthy(ctx context.Context) error {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := t.isHealthyRequest(ctx); err == nil {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (t *tester) isHealthyRequest(ctx context.Context) error {
	_, err := t.client.Check(ctx, &healther.HealthCheckRequest{})
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	return nil
}
