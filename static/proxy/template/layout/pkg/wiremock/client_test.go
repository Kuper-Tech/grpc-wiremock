package wiremock

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/require"
)

func Test_enrichWithMetaData(t *testing.T) {
	tests := []struct {
		name string
		md   metadata.MD

		wantHost   string
		wantHeader http.Header
	}{
		{
			name:       "empty metadata means empty host and header",
			md:         metadata.MD{},
			wantHost:   "",
			wantHeader: http.Header{},
		},
		{
			name: "authority metadata fills request host",
			md: metadata.MD{
				":authority": []string{
					"push-sender.test",
				},
			},
			wantHost:   "push-sender.test",
			wantHeader: http.Header{},
		},
		{
			name: "custom metadata fills request header",
			md: metadata.MD{
				"custom": []string{"testHeader"},
				"content-type": []string{
					"application/grpc",
				},
			},
			wantHost: "",
			wantHeader: http.Header{
				"Custom":       []string{"testHeader"},
				"Content-Type": []string{"application/grpc"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := DefaultRequest(
				metadata.NewIncomingContext(context.Background(), tt.md),
				"/Notify", bytes.NewReader([]byte{}))
			require.NoError(t, err)

			request = enrichWithMetaData(request)

			require.Equal(t, tt.wantHost, request.Host)
			require.Equal(t, tt.wantHeader, request.Header)
		})
	}
}
