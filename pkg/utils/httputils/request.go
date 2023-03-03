package httputils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/avast/retry-go/v4"
)

const defaultAttemptsValue = 7

func DoPost(ctx context.Context, client http.Client, url string, reader io.Reader) (int, error) {
	httpRequest, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return 0, fmt.Errorf("new request: %w", err)
	}

	httpRequest = httpRequest.WithContext(ctx)

	var httpResponse *http.Response

	err = retry.Do(func() error {
		httpResponse, err = client.Do(httpRequest)
		if err != nil {
			return fmt.Errorf("do: %w", err)
		}

		return nil
	}, retry.Attempts(defaultAttemptsValue))

	if err != nil {
		return 0, fmt.Errorf("retry: %w", err)
	}

	if httpResponse == nil {
		return 0, fmt.Errorf("response is empty")
	}

	defer func() {
		if err = httpResponse.Body.Close(); err != nil {
			log.Println("close body:", err)
		}
	}()

	return httpResponse.StatusCode, nil
}

func AssertStatus(expected, actual int) error {
	if actual != expected {
		return fmt.Errorf("expected status: %d, actual status: %d", expected, actual)
	}

	return nil
}
