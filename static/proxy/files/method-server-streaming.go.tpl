package {{ .PackageHeader }}

import (
	"bytes"
	"fmt"
	"net/http"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"grpc-proxy/pkg/wiremock"

	{{ range .GoPackages }}
	"{{ . }}"
	{{- end }}
)

func (p *Service) {{ .Method }}(in *{{ .MethodInPackage }}.{{ .MethodInName }}, stream {{ .MethodPackage }}.{{ .Service }}_{{ .Method }}Server) error {
	const url = "{{ .URL }}"
	const streamCursor = 1

	ctx := stream.Context()

	unmarshalAndSend := func(responseBody []byte) error {
		var protoResponse {{ .MethodOutPackage }}.{{ .MethodOutName }}
		if processErr := protojson.Unmarshal(responseBody, &protoResponse); processErr != nil {
			return status.Error(http.StatusBadGateway, fmt.Sprintf("marshal json object to proto: %v", processErr))
		}
		if processErr := stream.Send(&protoResponse); processErr != nil {
			return processErr
		}
		return nil
	}

	processStream := func(cursor int) error {
		httpRequest, processErr := wiremock.RequestWithCursor(ctx, url, cursor, bytes.NewReader([]byte{}))
		if processErr != nil {
			return status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", processErr))
		}
		httpResponseBody, processErr := wiremock.DoRequestDefault(httpRequest)
		if processErr != nil {
			return processErr
		}
		return unmarshalAndSend(httpResponseBody)
	}

	defaultRequest, err := wiremock.RequestWithCursor(ctx, url, streamCursor, bytes.NewReader([]byte{}))
	if err != nil {
		return status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", err))
	}

	httpResponseBody, streamSize, err := wiremock.DoRequestWithStreamSize(defaultRequest)
	if err != nil {
		return err
	}

	if err = unmarshalAndSend(httpResponseBody); err != nil {
		return err
	}

	for cursor := streamCursor + 1; cursor <= streamSize; cursor++ {
		if err = processStream(cursor); err != nil {
			return err
		}
	}

	return nil
}
