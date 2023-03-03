package {{ .PackageHeader }}

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"grpc-proxy/pkg/wiremock"

	{{ range .GoPackages }}
	"{{ . }}"
	{{- end }}
)

func (p *Service) {{ .Method }}(stream {{ .MethodPackage }}.{{ .Service }}_{{ .Method }}Server) error {
	const url = "{{ .URL }}"

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

	streamCursor := 1

	request, err := wiremock.RequestWithCursor(ctx, url, streamCursor, bytes.NewReader([]byte{}))
	if err != nil {
		return status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", err))
	}

	_, streamSize, err := wiremock.DoRequestWithStreamSize(request)
	if err != nil {
		return err
	}

	for {
		req, errReceive := stream.Recv()
		if errReceive != nil && errReceive == io.EOF {
			return nil
		}
		if errReceive != nil {
			return errReceive
		}
		if req == nil {
			continue
		}
		if err = processStream(streamCursor); err != nil {
			return err
		}
		if streamCursor >= streamSize {
			return nil
		}
		streamCursor++
	}
}
