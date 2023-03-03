package {{ .PackageHeader }}

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"grpc-proxy/pkg/wiremock"

	{{ range .GoPackages }}
	"{{ . }}"
	{{- end }}
)

func (p *Service) {{ .Method }}(ctx context.Context, in *{{ .MethodInPackage }}.{{ .MethodInName }}) (*{{ .MethodOutPackage }}.{{ .MethodOutName }}, error) {
	const url = "{{ .URL }}"

	requestBody, err := protojson.Marshal(in)
	if err != nil {
		return nil, status.Error(http.StatusBadGateway, fmt.Sprintf("create http request body: %v", err))
	}

	request, err := wiremock.DefaultRequest(ctx, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", err))
	}

	httpResponseBody, err := wiremock.DoRequestDefault(request)
	if err != nil {
		return nil, err
	}

	var protoResponse {{ .MethodOutPackage }}.{{ .MethodOutName }}
	if err = protojson.Unmarshal(httpResponseBody, &protoResponse); err != nil {
		return nil, status.Error(http.StatusBadGateway, fmt.Sprintf("marshal json object to proto: %v", err))
	}

	return &protoResponse, nil
}
