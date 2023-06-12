package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/mock"
)

func Test_newName(t *testing.T) {
	tests := []struct {
		name string

		mock mock.Mock

		want string
	}{
		{
			mock: mock.Mock{RequestMethod: "POST", RequestUrlPath: "/Example/RpcName-1", ResponseStatusCode: 200},
			want: "post_example_rpc_name_1_200.json",
		},
		{
			mock: mock.Mock{RequestMethod: "GET", RequestUrlPath: "api_rpc_name", ResponseStatusCode: 404},
			want: "get_api_rpc_name_404.json",
		},
		{
			mock: mock.Mock{RequestMethod: "GET", RequestUrlPath: "api_rpc_name/domain/handler", ResponseStatusCode: 404},
			want: "get_api_rpc_name_domain_handler_404.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createMockName(tt.mock)
			require.Equal(t, tt.want, got)
		})
	}
}
